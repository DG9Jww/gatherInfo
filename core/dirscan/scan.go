package dirscan

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/xuri/excelize/v2"
)

//Either dictionary or list must be set , but not both
func checkConfig(cfg *config.DirScanConfig) bool {
	logger.ConsoleLog(logger.NORMAL, "Checking DirScan Configuration......")
	if cfg.UrlList != nil && cfg.UrlDic != "" {
		logger.ConsoleLog(logger.ERROR, "Only one of urlDic and urlList is required.")
		return false
	}

	if cfg.PayloadDic == "" {
		logger.ConsoleLog(logger.ERROR, "DirScan payload dictionary must be set")
		return false
	}
	return true
}

//do task
func (cli *client) DoRequest(req *http.Request, wg *sync.WaitGroup, file *excelize.File, client *http.Client) func() {
	return func() {
		defer func() {
			wg.Done()
			atomic.AddInt64(&cli.done, 1)
		}()

		resp, err := client.Do(req)
		url := req.URL.String()
		//err != nil which means invalid
		if err != nil {

			//try http again
			url = strings.ReplaceAll(url, "https", "http")
			req, err := cli.GenerateRequest(url)
			if err != nil {
				return
			}
			resp, err = client.Do(req)
			if err != nil {
				return
			}
		}

		body := common.ReadHttpBody(resp)
		code := resp.StatusCode
		//if valid,processing the result
		var isValid bool
		if cli.validCode == nil {
			isValid = true
		} else {
			if common.MatchInt(code, cli.validCode) {
				isValid = true
			}
		}

		if isValid {

			//filter
			if len(cli.filterStr) > 0 {
				if common.MatchStr(cli.filterStr, string(body)) {
					return
				}
			}
			//30X process
			var r string
			if common.MatchInt(code, []int{301, 302, 303, 307}) {
				location, _ := resp.Location()
				r = location.String()
			}

			var res = &result{url: url, code: code, redirect: r, length: len(body)}
			cli.results <- res
		}

	}
}

func (cli *client) Scan(domain string, wg1 *sync.WaitGroup, client *http.Client) func() {
	return func() {
		var baseUrl string
		var signal = make(chan struct{})
		if strings.Contains(domain, "http") {
			baseUrl = domain
		} else {
			baseUrl = "https://" + domain
		}
		var wg sync.WaitGroup
		p := common.NewPool(cli.coroutine)

		//file for output
		var f *excelize.File
		if cli.output != "" {
			f = excelize.NewFile()
			f.SetSheetRow("Sheet1", "A1", &[]interface{}{"URL", "StatusCode", "Length", "Redirect"})
		}

		//a goroutine for processing results
		go func(file *excelize.File, signal chan struct{}) {
			var c int = 2
			for {
				select {
				case <-signal:
					return
				case res := <-cli.results:
					dirPrint(res)
					if cli.output != "" {
						file.SetSheetRow("Sheet1", fmt.Sprintf("A%d", c), &[]interface{}{res.url, res.code, res.length, res.redirect})
						c++
					}
				}
			}
		}(f, signal)

		for _, payload := range cli.payloadList {
			payload = url.QueryEscape(payload)
			url := baseUrl + "/" + payload
			req, err := cli.GenerateRequest(url)
			if err != nil {
				continue
			}
			p.Submit(cli.DoRequest(req, &wg, f, client))
			atomic.AddInt64(&cli.total, 1)
			wg.Add(1)
		}

		wg.Wait()
		//end processing results goroutine
		close(signal)
		if cli.output != "" {
			if _, err := os.Stat("output"); err != nil {
				os.Mkdir("output", 0755)
			}
			err := f.SaveAs("output/" + cli.output)
			if err != nil {
				logger.ConsoleLog(logger.ERROR, fmt.Sprintf("Save file %s failed :%s", cli.output, err.Error()))
			} else {
				logger.ConsoleLog(logger.INFO, fmt.Sprintf("Output file was save as %s ", cli.output))
			}
		}
		wg1.Done()
	}
}

func (cli *client) Run() {
	logger.ConsoleLog(logger.NORMAL, "Dirscan is Running......")

	if strings.EqualFold(cli.method, "POST") {
		logger.ConsoleLog(logger.ERROR, "POST method is not supported.")
	}
	//load payload file
	handle := common.LoadFile(cli.payloadDic)
	cli.payloadList = common.FileToSlice(handle)

	//This pool is for subdomains,I think 30 is enough
	var wg sync.WaitGroup
	p := common.NewPool(30)
	defer p.Release()

	//generate Client and set proxy
	client := common.NewHttpClient()
	if cli.proxy != "" {
		var proxy = func(*http.Request) (*url.URL, error) {
			return url.Parse(cli.proxy)
		}
		t := &http.Transport{Proxy: proxy, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client.Transport = t
	}

	//progress bar
	var endSignal = make(chan struct{})
	go func() {
		for {
			select {
			case <-endSignal:
				return
			default:
				time.Sleep(time.Millisecond * 50)
				fmt.Printf("\r[%d/%d] Done/Total", cli.done, cli.total)
			}
		}
	}()

	//if exist url dictionary
	if cli.urlDic != "" {
		handle := common.LoadFile(cli.urlDic)
		cli.urlList = common.FileToSlice(handle)
	}
	for _, v := range cli.urlList {
		p.Submit(cli.Scan(v, &wg, client))
		wg.Add(1)
	}

	wg.Wait()
	close(endSignal)
}

func Run(cfg *config.DirScanConfig, wg *sync.WaitGroup) {
	if cfg.Enabled {
		cli := NewClient(cfg)
		if checkConfig(cfg) {
			cli.Run()
		} else {
			logger.ConsoleLog(logger.WARN, "Please check your config or params !")
		}
	}
	wg.Done()
}

func (cli *client) ProcessResult() {}
