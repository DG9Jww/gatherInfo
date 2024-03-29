/*
CopyRight 2022
Author:DG9J
*/

package subdomain

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/subdomain/apis"
	"github.com/DG9Jww/gatherInfo/core/subdomain/enumerate"
	"github.com/DG9Jww/gatherInfo/core/subdomain/result"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/xuri/excelize/v2"
)

func Run(cfg *config.SubDomainConfig, isDir bool, wg *sync.WaitGroup) {
	//api result slice
	var subdomainList []string
	//for validate && output
	var resList []*result.Result
	//for duplicate
	var tmpSlice []string

	if cfg.Enabled {
		//process results
		go func() {
			for r := range result.FinalResults {
				//remove duplicate
				if common.IsStringInSlice(r.GetSubdomain(), tmpSlice) {
					continue
				}
				tmpSlice = append(tmpSlice, r.GetSubdomain())
				s := fmt.Sprintf("%s[%s]%s %s", logger.BLUE, r.GetSubdomain(), logger.ENDC, r.GetRecord())
				logger.ConsoleLog(logger.CustomizeLog(logger.WHITE, fmt.Sprintf("\r[%s%s%s]", logger.WHITE, "+", logger.ENDC)), s)
				resList = append(resList, r)
			}
		}()

		//different mode
		switch cfg.Mode {
		case "":
			subdomainList = apis.Run(cfg.Domain)
			for _, subdomain := range subdomainList {
				ip, _ := net.LookupIP(subdomain)
				for _, v := range ip {
					if v.To4() != nil {
						var r = &result.Result{}
						r.SetSubdomain(subdomain)
						r.SetRecord(v.To4().String())
						result.FinalResults <- r
					}
				}
			}
			enumerate.Run(cfg)

		case "api":
			subdomainList = apis.Run(cfg.Domain)
			for _, subdomain := range subdomainList {
				ip, _ := net.LookupIP(subdomain)
				for _, v := range ip {
					if v.To4() != nil {
						var r = &result.Result{}
						r.SetSubdomain(subdomain)
						r.SetRecord(v.To4().String())
						result.FinalResults <- r
					}
				}
			}

		case "enu":
			enumerate.Run(cfg)
		}
		close(result.FinalResults)

		var total = len(tmpSlice)
		logger.ConsoleLog(logger.INFO, fmt.Sprintf("===== %d Subdomain Found =====", total))

		//validate whether the subdomain is live
		if cfg.Validate {
			var done int
			var barSignal = make(chan struct{})
			//progress bar
			go func(signal chan struct{}) {
				for {
					select {
					case <-signal:
						return
					default:
						time.Sleep(time.Millisecond * 50)
						fmt.Printf("\r[%d/%d]", done, total)
					}
				}
			}(barSignal)
			pool := common.NewPool(20)
			defer pool.Release()
			var wg sync.WaitGroup
			for _, v := range resList {
				pool.Submit(isLive(v, &wg, &done))
				wg.Add(1)
			}
			wg.Wait()
			close(barSignal)
		}

		//out put excel file
		if cfg.OutPut != "" {
			f := excelize.NewFile()
			f.SetSheetRow("Sheet1", "A1", &[]interface{}{"SubdomainName", "ParsedResult", "StatusCode"})
			for index, res := range resList {
				f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", index+2), &[]interface{}{res.GetSubdomain(), res.GetRecord(), res.GetCode()})
			}
			if _, err := os.Stat("output"); err != nil {
				os.Mkdir("output", 0755)
			}
			if err := f.SaveAs("output/" + cfg.OutPut); err != nil {
				logger.ConsoleLog(logger.ERROR, fmt.Sprintf("Save file %s error:%s", cfg.OutPut, err.Error()))
			} else {
				logger.ConsoleLog(logger.INFO, fmt.Sprintf("Output file was save as %s", cfg.OutPut))
			}
		}
	}
	wg.Done()
}

// validate whether the url is live
func isLive(res *result.Result, wg *sync.WaitGroup, done *int) func() {
	return func() {
		defer func() {
			wg.Done()
			*done++
		}()
		subdomain := res.GetSubdomain()
		url := fmt.Sprintf("https://%s", subdomain)
		req, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		var code int
		var resp *http.Response
		resp, err = common.HttpRequest(req)
		if err != nil {
			url = fmt.Sprintf("http://%s", subdomain)
			req2, err := common.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			resp, err = common.HttpRequest(req2)
			if err != nil {
				return
			}
		}
		code = resp.StatusCode
		logger.StatusCodeLog(code, url)
		res.SetCode(code)
		return
	}
}
