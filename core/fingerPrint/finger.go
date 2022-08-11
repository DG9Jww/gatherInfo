package fingerprint

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
)

type fingerPrint struct {
	Path    string
	Option  string //md5 and keyword
	Content string
	Name    string
}

func Run(cfg *config.FingerPrintConfig, wg *sync.WaitGroup) {
	if cfg.Enabled {
		cli := NewClient(cfg)
		cli.Run()
		logger.ConsoleLog(logger.NORMAL, "FINGERPRINT COMPLETED")
	}
	wg.Done()
}

func (cli *client) Run() {
	if err := cli.loadFingerPrint(); err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	logger.ConsoleLog(logger.NORMAL, "FingerPrint Module is Runnning......")
	pool := common.NewPool(cli.thread)
	defer pool.Release()

	//let's go
	urlList := cli.getUrlList()
	for _, url := range urlList {
		for cmsName, fingerPList := range cli.fingerMap {
			for _, fingerP := range fingerPList {
				fingerP.Name = cmsName

				tempURL := url + fingerP.Path
				cli.wg.Add(1)
				pool.Submit(cli.Scan(tempURL, fingerP))
			}

		}
	}
	cli.wg.Wait()
}

func (cli *client) Scan(url string, f fingerPrint) func() {
	return func() {
		defer cli.wg.Done()
		var resp *http.Response
		req, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		//default https
		resp, err = common.HttpRequest(req)
		if err != nil {

			//try http
			url = strings.ReplaceAll(url, "https://", "http://")
			req, err := common.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			resp, err = common.HttpRequest(req)
			if err != nil {
				return
			}
		}

		//response code require 200
		if resp.StatusCode != 200 {
			return
		}

		data := common.ReadHttpBody(resp)
		//match fingerprint
		option := f.Option
		if option == "keyword" {
			if common.MatchStr(f.Content, string(data)) {
				logger.ConsoleLog(logger.FINGERPRINT, fmt.Sprintf("[FINGERPRINT]  CMS:%s  URL:%s", f.Name, url))
			}
		} else if option == "md5" {
			md5Str := fmt.Sprintf("%x", md5.Sum(data))
			if md5Str == f.Content {
				logger.ConsoleLog(logger.FINGERPRINT, fmt.Sprintf("[FINGERPRINT]  CMS:%s  URL:%s", f.Name, url))
			}
		}
	}
}

//parse fingerprint  into struct
func (cli *client) loadFingerPrint() error {
	data, err := os.ReadFile(cli.fingerP)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &cli.fingerMap)
	if err != nil {
		return nil
	}
	return err
}

func (cli *client) getUrlList() (urlList []string) {
	for _, url := range cli.urlList {
		var u string
		if strings.Contains(url, "http") {
			u = strings.ReplaceAll(url, "http://", "https://")

		} else {
			u = "https://" + url
		}
		urlList = append(urlList, u)
	}
	return urlList
}

func (cli *client) domainToUrl(domain string) (url string) {
	u := "https://" + domain
	return u
}
