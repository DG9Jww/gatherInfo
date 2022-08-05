package apis

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/DG9Jww/gatherInfo/logger"
)

//function
const (
	funcBase64 = `$base64`
	rootDir    = `core/subdomain/apis/scripts`
)

var (
	funcList = make(map[string]func(string) string)
)

type APIRequest struct {
	BaseUrl   string                 `json:"baseurl"`
	Path      string                 `json:"path"`
	Method    string                 `json:"method"`
	Headers   map[string]string      `json:"headers"`
	Variables map[string]string      `json:"variables"`
	PostBody  map[string]interface{} `json:"postbody"`
}

//json process and send request
func start(APIName string, data []byte, domain string, wg *sync.WaitGroup) {
	defer wg.Done()
	var req = new(APIRequest)
	err := json.Unmarshal(data, req)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	//looking for function && variables and process
	//Firstly,we must add domain name into the map
	req.Variables["domain"] = domain
	data1 := req.replaceVariables(data)
	data2 := req.runFunc(funcList, data1)

	//marshal json data again
	err = json.Unmarshal([]byte(data2), req)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	//send request
	resp, err := req.sendRequest()
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}
	defer resp.Body.Close()

	//process response
	processResp(APIName, resp)

	//fmt.Println(string(b))
}

//looking for and replace variables
func (req *APIRequest) replaceVariables(data []byte) string {
	tmpData := string(data)
	for k, v := range req.Variables {
		tmpData = regexp.MustCompile(fmt.Sprintf("{{%s}}", k)).ReplaceAllString(string(tmpData), v)
	}
	return tmpData
}

//looking for and run function
func (req *APIRequest) runFunc(funcList map[string]func(string) string, str string) string {
	for funcStr, f := range funcList {
		out := regexp.MustCompile(fmt.Sprintf("\\%s\\((.*?)\\)", funcStr)).FindAllStringSubmatch(str, -1)
		if out == nil {
			continue
		}
		for _, v := range out {
			//the string in the braces
			tmp := v[1]
			funcOut := f(tmp)
			str = regexp.MustCompile(fmt.Sprintf("\\%s\\((%s)\\)", funcStr, tmp)).ReplaceAllString(str, funcOut)
		}
	}
	return str
}
