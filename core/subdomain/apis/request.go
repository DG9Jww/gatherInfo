package apis

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

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
	err := json.Unmarshal(b, apiStruct)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	//return the value that the pointer points to
	v := reflect.Indirect(reflect.ValueOf(apiStruct))

	/*
	 *
	 *    --------------------------------------------
	 *       ptr
	 *       ↓
	 *      -----            -----
	 *      reflect*Value
	 *      -----            -----
	 *      valChan          valChan      valChan x n
	 *    --------------------------------------------
	 *           structChan
	 */
	var maximum = 15
	type valChan chan reflect.Value
	var structChan = make(chan valChan, 10)
	var vChan valChan = make(chan reflect.Value, maximum)
	structChan <- vChan
	vChan <- v
	var ptr *valChan

	//for closing Chan
	go func() {
		for {
			if len(resSlice) > 0 {
				close(structChan)
				close(*ptr)
				return
			}
		}
	}()

	for vChan := range structChan {
		ptr = &vChan
		for v := range *ptr {

			// struct's fields
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
                
				//string (result)
				if field.Kind() == reflect.String {
					var s = Result{domain: field.String()}
					resSlice = append(resSlice, s)
				}

				//slice
				if field.Kind() == reflect.Slice {
					//[]string (result)
					if field.Type().Elem().Kind() == reflect.String {
						// slice's length
						l := field.Len()
						for j := 0; j < l; j++ {
							item := field.Index(j)
							var s = Result{domain: item.String()}
							resSlice = append(resSlice, s)
						}
					}

					//[]struct
					if field.Type().Elem().Kind() == reflect.Struct {
						// slice's length
						l := field.Len()
						for j := 0; j < l; j++ {
							item := field.Index(j)
							if len(*ptr) >= maximum {
								close(*ptr)
								var newChan valChan = make(chan reflect.Value, maximum)
								ptr = &newChan
								structChan <- newChan
								*ptr <- item
							} else {
								*ptr <- item
							}
						}
					}
				}
			}

		}

	}

}
