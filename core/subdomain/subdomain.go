/*
CopyRight 2022
Author:DG9J
*/

package subdomain

import (
	"sync"

	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/subdomain/enumerate"
	"github.com/DG9Jww/gatherInfo/core/subdomain/results"
)

var (
	SubDomainRes = results.NewResults()
)

func Run(cfg *config.SubDomainConfig, isDir bool, wg *sync.WaitGroup) {
	if cfg.Enabled {
		//logger.ConsoleLog(logger.NORMAL, "Subdomain module is Running......")
		////fofa
		//err := fofa.NewClient(SubDomainRes, cfg.FofaKey, cfg.FofaEmail, cfg.Domain).GetAPIInfo()
		//if err != nil {
		//	logger.ConsoleLog(logger.ERROR, err.Error())
		//	return
		//}

		////cert
		////	err = cert.NewClient(SubDomainRes, cfg.Domain, cfg.CensysID, cfg.CensysKey).Run()
		////	if err != nil {
		////		fmt.Println("11111111111")
		////		logger.ConsoleLog(logger.ERROR, err.Error())
		////		fmt.Println("2222222222")
		////		return
		////	}

		//SubDomainRes.RemoveDuplicate()
		//SubDomainRes.VerifyDomain(isDir)

		//brute module
		enumerate.Run(cfg)
	}
	wg.Done()
}
