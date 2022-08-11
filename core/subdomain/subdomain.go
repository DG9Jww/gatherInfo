/*
CopyRight 2022
Author:DG9J
*/

package subdomain

import (
	"fmt"
	"net"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/subdomain/apis"
	"github.com/DG9Jww/gatherInfo/core/subdomain/enumerate"
	"github.com/DG9Jww/gatherInfo/logger"
)

func Run(cfg *config.SubDomainConfig, isDir bool, wg *sync.WaitGroup) {
	var subdomainList []string
	if cfg.Enabled {
		switch cfg.Mode {
		case "":
			//apis
			subdomainList = apis.Run(cfg.Domain)
			//brute module
			enumerate.Run(cfg)
		case "api":
			subdomainList = apis.Run(cfg.Domain)
		case "enu":
			enumerate.Run(cfg)
		}
	}

	for _, subdomain := range subdomainList {
		logger.ConsoleLog2(logger.CustomizeLog(logger.BLUE, subdomain), "")
	}
	for _, subdomain := range subdomainList {

		ip, _ := net.LookupIP(subdomain)
		for _, v := range ip {
			fmt.Println(v.To4())
		}

	}

    //validate whether the subdomain is live
    if cfg.Validate {

    }

	wg.Done()
}


// validate whether the url is live
func isLive(subdomain string, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		url := fmt.Sprintf("https://%s", subdomain)
		req, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		resp, err := common.HttpRequest(req)
		if err != nil {
			url2 := fmt.Sprintf("http://%s", subdomain)
			req2, err := common.NewRequest("GET", url2, nil)
			if err != nil {
				return
			}
			resp2, err := common.HttpRequest(req2)
			if err != nil {
				return
			}
			logger.StatusCodeLog(resp2.StatusCode, url2)
			return
		}
		logger.StatusCodeLog(resp.StatusCode, url)
		return
	}
}
