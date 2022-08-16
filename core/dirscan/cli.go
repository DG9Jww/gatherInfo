package dirscan

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
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
	results chan *result

	//count completed task
	done int64

	//count completed task
	total int64

	//valid statuscode
	validCode []int

	//filter string
	filterStr string

	//proxy
	proxy string

	//method
	method string

	//header
	header string

	//output
	output string
}

type result struct {
	url      string
	code     int
	length   int
	redirect string
}

func NewClient(cfg *config.DirScanConfig) *client {
	c := &client{
		results:     make(chan *result),
		coroutine:   cfg.Coroutine,
		urlDic:      cfg.UrlDic,
		urlList:     cfg.UrlList,
		payloadDic:  cfg.PayloadDic,
		payloadList: cfg.PayloadList,
		validCode:   cfg.ValidCode,
		filterStr:   cfg.FilterStr,
		proxy:       cfg.Proxy,
		method:      cfg.Method,
		header:      cfg.Header,
		output:      cfg.OutPut,
	}
	return c
}

func (cli *client) GenerateRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(cli.method, url, nil)
	if err != nil {
		return nil, err
	}
	if cli.header != "" {
		tmp := strings.Split(cli.header, ":")
		key := strings.TrimSpace(tmp[0])
		value := strings.TrimSpace(tmp[1])
		req.Header.Add(key, value)
	} else {
		req.Header.Set("User-Agent", common.RandomAgent())
	}
	if err != nil {
		return nil, err
	}
	return req, err
}

func dirPrint(res *result) {
	switch {
	case res.code >= 200 && res.code < 300:
		logger.ConsoleLog(logger.R20X, fmt.Sprintf("\r[+] %d %s  Length:%d", res.code, res.url, res.length))
	case res.code >= 300 && res.code < 400:
		logger.ConsoleLog(logger.R30X, fmt.Sprintf("\r[+] %d %s ==> %s Length:%d", res.code, res.url, res.redirect, res.length))
	case res.code >= 400 && res.code < 500:
		logger.ConsoleLog(logger.R40X, fmt.Sprintf("\r[+] %d %s  Length:%d", res.code, res.url, res.length))
	case res.code >= 500:
		logger.ConsoleLog(logger.R50X, fmt.Sprintf("\r[+] %d %s  Length:%d", res.code, res.url, res.length))
	}
}
