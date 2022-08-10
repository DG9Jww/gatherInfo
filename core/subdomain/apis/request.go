package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DG9Jww/gatherInfo/common"
)

type APIRequest struct {
	BaseUrl      string                 `json:"baseurl"`
	Path         string                 `json:"path"`
	Method       string                 `json:"method"`
	Headers      map[string]string      `json:"headers"`
	Variables    map[string]string      `json:"variables"`
	PostBody     map[string]interface{} `json:"postbody"`
	NeedRE       ReField                `json:"needre"`
	ResponseType string                 `json:"response_type"`
}

type ReField struct {
	Subdomain bool `json:"subdomain"`
	IP        bool `json:"ip"`
}

func (req *APIRequest) sendRequest() (*http.Response, error) {
	url := req.BaseUrl + req.Path
    fmt.Println(url)
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

