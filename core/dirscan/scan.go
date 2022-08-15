package dirscan

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
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
func (cli *client) DoRequest(req *http.Request, wg *sync.WaitGroup, file *os.File, client *http.Client) func() {
	return func() {
		defer func() {
			wg.Done()
			cli.lock.Lock()
			cli.counter++
			cli.lock.Unlock()
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
		}
		if common.MatchInt(code, cli.validCode) {
			isValid = true
		}

		if isValid {
			temp := fmt.Sprintf("%d %s", code, url)

			//filter
			if len(cli.filterStr) > 0 {
				if common.MatchStr(cli.filterStr, string(body)) {
					return
				}
			}
			//30X process
			if common.MatchInt(code, []int{301, 302, 303, 307}) {
				location, _ := resp.Location()
				l := location.String()
				temp = fmt.Sprintf("%d %s ===> %s", code, url, l)
			} else {
				temp = fmt.Sprintf("%d %s Length: %d", code, url, len(body))
			}

			cli.lock.Lock()
			cli.results = append(cli.results, temp)
			cli.lock.Unlock()
			dirPrint(code, file, temp)

		}

	}
}

func (cli *client) Scan(domain string, wg1 *sync.WaitGroup, client *http.Client) func() {
	return func() {
		var baseUrl string
		if strings.Contains(domain, "http") {
			baseUrl = domain
		} else {
			baseUrl = "https://" + domain
		}
		var wg sync.WaitGroup
		p := common.NewPool(cli.coroutine)

		//file for output
		path := strings.Replace(domain, "://", "_", -1)
		path = "output/" + path + ".txt"
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logger.ConsoleLog(logger.ERROR, err.Error())
		}
		for _, payload := range cli.payloadList {
			url := baseUrl + "/" + payload
			req, err := cli.GenerateRequest(url)
			if err != nil {
				continue
			}
			p.Submit(cli.DoRequest(req, &wg, file, client))
			wg.Add(1)
		}

		wg.Wait()
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
			return url.Parse("http://127.0.0.1:8080")
		}
		t := &http.Transport{Proxy: proxy, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client.Transport = t
	}

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
