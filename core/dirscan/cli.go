package dirscan

import (
	"os"
	"sync"

	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
)

type client struct {
	lock sync.Mutex

	//for path scan
	coroutine int

	//target dictionary
	urlDic string

	//target list
	urlList []string

	//payload list for scan
	payloadDic string

	//payload list
	payloadList []string

	//results
	results []string

	//count completed task
	counter int

	//valid statuscode
	validCode []int

	//filter string
	filterStr string
}

func NewClient(cfg *config.DirScanConfig) *client {
	c := &client{
		coroutine:   cfg.Coroutine,
		urlDic:      cfg.UrlDic,
		urlList:     cfg.UrlList,
		payloadDic:  cfg.PayloadDic,
		payloadList: cfg.PayloadList,
		validCode:   cfg.ValidCode,
		filterStr:   cfg.FilterStr,
	}
	return c
}

func dirPrint(code int, file *os.File, content string) {
	switch {
	case code >= 200 && code < 300:
		logger.ConsoleLog(logger.R20X, content)
	case code >= 300 && code < 400:
		logger.ConsoleLog(logger.R30X, content)
	case code >= 400 && code < 500:
		logger.ConsoleLog(logger.R40X, content)
	case code >= 500:
		logger.ConsoleLog(logger.R50X, content)
	}
	logger.LogToFile(file, content)
}
