package fingerprint

import (
	"sync"

	"github.com/DG9Jww/gatherInfo/config"
)

type client struct {
	thread int

	//dictionary path
	fingerP string

	//fingerprint struct
	fingerMap map[string][]fingerPrint

	//target urls
	urlList []string

	wg      sync.WaitGroup
	results []string
}

func NewClient(cfg *config.FingerPrintConfig) *client {
	cli := &client{
		thread:  cfg.Thread,
		fingerP: cfg.FingerP,
		urlList: cfg.UrlList,
	}
	return cli
}
