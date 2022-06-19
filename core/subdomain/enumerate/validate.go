package enumerate

import (
	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
)

// validate whether the url is live
func isLive(url string) {
	req, err := common.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	resp, err := common.HttpRequest(req)
	if err != nil {
		return
	}
	logger.ConsoleLog(logger.NORMAL, resp.Status)
}
