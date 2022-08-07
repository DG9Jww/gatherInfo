package apis

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/DG9Jww/gatherInfo/common"
)

func (req *APIRequest) sendRequest() (*http.Response, error) {
	url := req.BaseUrl + req.Path
	var apiReq *http.Request
	var err error
	if len(req.PostBody) > 0 {
		d, err := json.Marshal(req.PostBody)
		if err != nil {
			return nil, err
		}
		apiReq, err = http.NewRequest(req.Method, url, bytes.NewReader(d))
		if err != nil {
			return nil, err
		}
	} else if len(req.PostBody) == 0 {
		apiReq, err = http.NewRequest(req.Method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	//add Headers
	if req.Headers != nil {
		for k, v := range req.Headers {
			apiReq.Header.Add(k, v)
		}
	}

	return common.HttpRequest(apiReq)
}

