package dirscan

import (
	"os"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/subdomain"
	"github.com/DG9Jww/gatherInfo/logger"
)

//Either dictionary or list must be set , but not both
func checkConfig(cfg *config.DirScanConfig, isSubDomain bool) bool {
	logger.ConsoleLog(logger.NORMAL, "Checking DirScan Configuration......")
	if cfg.UrlList != nil && cfg.UrlDic != "" {
		logger.ConsoleLog(logger.ERROR, "Only one of urlDic and urlList is required.")
		return false
	}

	//When you using dirscan module and did't use subdomain module,,
	//urlList or urlDic is required.
	if isSubDomain {
		if cfg.UrlList == nil && cfg.UrlDic == "" {
			logger.ConsoleLog(logger.ERROR, "UrlDic or UrlList is required.")
			return false
		}
	}

	if cfg.PayloadDic == "" {
		logger.ConsoleLog(logger.ERROR, "DirScan payload dictionary must be set")
		return false
	}
	return true
}

//do task
func (cli *client) DoRequest(url string, wg *sync.WaitGroup, file *os.File) func() {
	return func() {
		defer func() {
			wg.Done()
			cli.lock.Lock()
			cli.counter++
			cli.lock.Unlock()
		}()

		r, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		resp, err := common.HttpRequest(r)

		//err != nil which means invalid
		if err != nil {

			//try http again
			url = strings.ReplaceAll(url, "https", "http")
			r, err := common.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			resp, err = common.HttpRequest(r)
			if err != nil {
				return
			}
		}

		body := common.ReadHttpBody(resp)
		code := resp.StatusCode
		//if valid,processing the result
		if common.MatchInt(code, cli.validCode) {
			temp := []interface{}{code, url}
			//filter
			if len(cli.filterStr) > 0 {
				if common.MatchStr(cli.filterStr, string(body)) {
					return
				}
			}
			//30X process
			if common.MatchInt(code, []int{301, 302, 303, 307}) {
				location, _ := resp.Location()
				temp = []interface{}{code, url, "==>", location}
			} else {
				temp = []interface{}{code, url, "Length:", len(body)}
			}

			cli.lock.Lock()
			cli.results = append(cli.results, temp)
			cli.lock.Unlock()
			dirPrint(code, file, temp...)

		}

	}
}

func (cli *client) Scan(domain string, wg1 *sync.WaitGroup) func() {
	return func() {
		var baseUrl string
		if strings.Contains(domain, "http") {
			baseUrl = domain
		} else {
			baseUrl = "https://" + domain
		}
		var wg sync.WaitGroup
		//Every subdomain is a single goroutine
		//and every subdomain goroutine need multiple sub-goroutine
		//to run "DirScan".Just like a nesting goroutine
		p := common.NewPool(cli.coroutine)

		//file for output
		path := strings.Replace(domain, "://", "_", -1)
		path = "output/" + path + ".txt"
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logger.ConsoleLog(logger.ERROR, err)
		}
		for _, payload := range cli.payloadList {
			url := baseUrl + "/" + payload
			p.Submit(cli.DoRequest(url, &wg, file))
			wg.Add(1)
		}

		wg.Wait()
		wg1.Done()
	}
}

func (cli *client) Run(isSubDomain bool) {
	logger.ConsoleLog(logger.NORMAL, "Dirscan is Running......")

	//load payload file
	handle := common.LoadFile(cli.payloadDic)
	cli.payloadList = common.FileToSlice(handle)

	//This pool is for subdomains,I think 30 is enough
	var wg sync.WaitGroup
	p := common.NewPool(30)
	defer p.Release()

	//If subdomain module has already been enable,domainList in dirScan
	//will be replaced by subdomain module results
	if isSubDomain {
		dirChan := subdomain.SubDomainRes.GetDomainChan()
		for domain := range dirChan {
			p.Submit(cli.Scan(domain, &wg))
			wg.Add(1)
		}
	} else {
		logger.ConsoleLog(logger.WARN, "subdomain did not enable")
		//if exist url dictionary
		if cli.urlDic != "" {
			handle := common.LoadFile(cli.urlDic)
			cli.urlList = common.FileToSlice(handle)
		}
		for _, v := range cli.urlList {
			p.Submit(cli.Scan(v, &wg))
			wg.Add(1)
		}
	}

	wg.Wait()
}

func Run(cfg *config.DirScanConfig, isSubDomain bool, wg *sync.WaitGroup) {
	if cfg.Enabled {
		cli := NewClient(cfg)
		if checkConfig(cfg, isSubDomain) {
			cli.Run(isSubDomain)
		} else {
			logger.ConsoleLog(logger.WARN, "Please check your config or params !")
		}
	}
	wg.Done()
}
