/*
CopyRight 2022
Author:DG9J
*/
package core

import (
	"os"
	"sync"

	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/dirscan"
	fingerprint "github.com/DG9Jww/gatherInfo/core/fingerPrint"
	"github.com/DG9Jww/gatherInfo/core/portscan"
	"github.com/DG9Jww/gatherInfo/core/subdomain"
	"github.com/DG9Jww/gatherInfo/logger"
)

func Execute() {
	var cfg *config.MyConfig
	//configuration file mode
	if len(os.Args) == 1 {
		logger.ConsoleLog(logger.NORMAL, "Using Configuration File Mode")
		cfg = config.ConfigFileInit()
		config.Mode = 1
	} else if len(os.Args) > 1 {
		//command mode
		logger.ConsoleLog(logger.NORMAL, "Using Command Mode")
		cfg = config.ConfigCommandInit()
		config.Mode = 0
	}

	//Run Modules
	var wg sync.WaitGroup
	wg.Add(4)
	go subdomain.Run(&cfg.SubDomain, cfg.DirScan.Enabled, &wg)
	go dirscan.Run(&cfg.DirScan, &wg)
	go portscan.Run(&cfg.PortScan, &wg)
	go fingerprint.Run(&cfg.FingerPrint, &wg)
	wg.Wait()
}
