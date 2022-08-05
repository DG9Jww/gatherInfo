package apis

import (
	"bytes"
    "io"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
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

//process response
func processResp(APIName string, resp *http.Response) {

	b, _ := io.ReadAll(resp.Body)
	apiStruct := APIStruct[APIName]
	fmt.Printf("%T\n", apiStruct)
	xxx(apiStruct)
	fmt.Printf("%T\n", apiStruct)
	//fmt.Println(reflect.ValueOf(apiStruct))
	//instance := reflect.New(apiStruct).Elem()
	//tmpStruct := instance.Interface()

    var oo = virustotal{}
    
	err := json.Unmarshal(b, &oo)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}
	fmt.Println(oo)

	//v := reflect.ValueOf(tmpStruct)
	//var structChan = make(chan reflect.Value, 10)
	//structChan <- v

	////for closing  structChan
	//go func() {
	//	for {
	//		if len(resSlice) > 0 {
	//			close(structChan)
	//			return
	//		}
	//	}
	//}()

	//for v := range structChan {
	//	for i := 0; i < v.NumField(); i++ {
	//		field := v.Field(i)

	//		//string (result)
	//		if field.Kind() == reflect.String {
	//			var s = Result{domain: v.String()}
	//			resSlice = append(resSlice, s)
	//		}

	//		//slice
	//		if field.Kind() == reflect.Slice {
	//			//[]string (result)
	//			if field.Type().Elem().Kind() == reflect.String {
	//				// slice's length
	//				l := field.Len()
	//				for j := 0; j < l; j++ {
	//					item := field.Index(j)
	//					var s = Result{domain: item.String()}
	//					resSlice = append(resSlice, s)
	//				}
	//			}

	//			//[]struct
	//			if field.Type().Elem().Kind() == reflect.Struct {
	//				// slice's length
	//				l := field.Len()
	//				for j := 0; j < l; j++ {
	//					item := field.Index(j)
	//					structChan <- item
	//				}
	//			}
	//		}
	//	}
	//}

	//fmt.Println(resSlice)

}
