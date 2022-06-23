package enumerate

import (
	"fmt"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
)

// validate whether the url is live
func isLive(subdomain string, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		url := fmt.Sprintf("https://%s", subdomain)
		req, err := common.NewRequest("HEAD", url, nil)
		if err != nil {
			return
		}
		resp, err := common.HttpRequest(req)
		if err != nil {
			url2 := fmt.Sprintf("http://%s", subdomain)
			req2, err := common.NewRequest("HEAD", url2, nil)
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
