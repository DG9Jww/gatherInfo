package apis

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/logger"
)




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
	if req.Variables == nil {
		req.Variables = make(map[string]string)
	}
	req.Variables["domain"] = domain
	data1 := req.replaceVariables(data)
    //escape
    s := escape(data1)
	data2 := req.runFunc(funcList, s)

	//marshal json data again
	err = json.Unmarshal([]byte(data2), req)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	//send request
	resp, err := req.sendRequest()
	if err != nil {
		logger.ConsoleLog(logger.ERROR, fmt.Sprintf("API %s ERROR:%s", APIName, err.Error()))
		return
	}

	defer resp.Body.Close()

	//process response
	req.processResp(APIName, resp, domain)
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
			s, _ := strconv.Unquote("`" + tmp + "`")
			funcOut := f(s)
			str = regexp.MustCompile(fmt.Sprintf("\\%s\\((%s)\\)", funcStr, tmp)).ReplaceAllString(str, funcOut)
		}
	}
	return str
}

//escape
func escape(s string) string {
    return strings.ReplaceAll(s,`\"`,`"`) 
}
