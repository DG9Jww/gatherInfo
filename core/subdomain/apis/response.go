package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/DG9Jww/gatherInfo/logger"
)

//regular expression
const (
	getDomain = `[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z]{0,62})+\.?`
	getIP     = `[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}`
)

//process response
func (req *APIRequest) processResp(APIName string, resp *http.Response, domain string) {
	b, _ := io.ReadAll(resp.Body)
	switch req.ResponseType {
	case "json":
		proJsonResp(b, APIName, req.NeedRE)
	case "raw":
		proRawResp(b, domain,APIName)
	case "special":
		proSpecialResp(b, APIName, domain, req.NeedRE.Subdomain)
	}

    
}

//process response with particular function
func proSpecialResp(b []byte, APIName string, domain string, needRE bool) {
	p := SpecialRespMap[APIName]
	tmp, err := p.SpecialProcess(b)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}

	if needRE {
		exp := getExp(domain)
		tmp = proRegularExp(&tmp, exp)
	}
	for _, subdomain := range tmp {
		var res = Result{}
		res.domain = subdomain
        addResSlice(res)
	}
}

//process response without any format
func proRawResp(b []byte, domain string,name string) {

	resp := string(b)
	var tmpRes []string
	tmpRes = append(tmpRes, resp)
	//only match subdomain
	exp := getExp(domain)
	s := proRegularExp(&tmpRes, exp)
    fmt.Println("88888888",name,s)
	for _, subdomain := range s {
		var res = Result{}
		res.domain = subdomain
        addResSlice(res)
	}

}

//process json response
func proJsonResp(b []byte, APIName string, needRE ReField) {

	apiStruct := APIStruct[APIName]
	err := json.Unmarshal(b, apiStruct)
	if err != nil {
		logger.ConsoleLog(logger.WARN, fmt.Sprintf("API %s ERROR:%s", APIName, err.Error()))
		return
	}

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
	//return the value that the pointer points to
	v := reflect.Indirect(reflect.ValueOf(apiStruct))
	if v.IsZero() {
		logger.ConsoleLog(logger.WARN, fmt.Sprintf("The unmarshal results of API %s is nil", APIName))
		return
	}
	var maximum = 15
	type valChan chan reflect.Value
	var structChan = make(chan valChan, 10)
	var vChan valChan = make(chan reflect.Value, maximum)
	structChan <- vChan
	vChan <- v
	var ptr *valChan
	var subdomainSlice []string
	var ipaddressSlice []string

	//for closing Chan
	go func() {
		for {
			if len(*ptr) == 0 && (len(subdomainSlice) > 0 || len(ipaddressSlice) > 0) {
				close(structChan)
				close(*ptr)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	for vChan := range structChan {
		ptr = &vChan
		for v := range vChan {

			//struct
			if v.Kind() == reflect.Struct {
				// struct's fields
				for i := 0; i < v.NumField(); i++ {
					field := v.Field(i)
					fieldName := v.Type().Field(i).Name

					//string
					if field.Kind() == reflect.String {
						// result
						switch fieldName {
						case "Subdomain":
							subdomainSlice = append(subdomainSlice, field.String())
						case "IPaddress":
							ipaddressSlice = append(ipaddressSlice, field.String())
						case "DomainAndIP":
							subdomainSlice = append(ipaddressSlice, field.String())
							ipaddressSlice = append(ipaddressSlice, field.String())
						}
					}

					//slice
					if field.Kind() == reflect.Slice {

						//[]string
						if field.Type().Elem().Kind() == reflect.String {
							var tmpSlice *[]string
							switch fieldName {
							case "SubdomainSlice":
								tmpSlice = &subdomainSlice
							case "IPaddressSlice":
								tmpSlice = &ipaddressSlice
							}

							// slice's length
							l := field.Len()
							for j := 0; j < l; j++ {
								item := field.Index(j)
								*tmpSlice = append(*tmpSlice, item.String())
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
									structChan <- *ptr
									*ptr <- item
								} else {
									*ptr <- item
								}
							}
						}

					}

					//struct
					if field.Kind() == reflect.Struct {
						if len(*ptr) >= maximum {
							close(*ptr)
							var newChan valChan = make(chan reflect.Value, maximum)
							ptr = &newChan
							structChan <- *ptr
							*ptr <- field
						} else {
							*ptr <- field
						}
					}

				}
			}

		}

	}

	//regular expression process
	if needRE.IP || needRE.Subdomain {
		var (
			do       = true
			tmpSlice []*[]string
			exp      []string
		)

		switch do {
		case needRE.Subdomain:
			tmpSlice = append(tmpSlice, &subdomainSlice)
			exp = append(exp, getDomain)
		case needRE.IP:
			tmpSlice = append(tmpSlice, &ipaddressSlice)
			exp = append(exp, getIP)
		}

		for k, sliceAddr := range tmpSlice {
			*sliceAddr = proRegularExp(sliceAddr, exp[k])
		}
	}

	//append into reSlice
	for index, subdomain := range subdomainSlice {
		var res = Result{}
		res.domain = subdomain
		if len(ipaddressSlice) > 0 && index < len(ipaddressSlice) {
			res.ip = ipaddressSlice[index]
		}
        addResSlice(res)
	}
}

//process ip and domain according to regular expression
func proRegularExp(tmpResSlice *[]string, exp string) []string {
	var tmp []string
	re, err := regexp.Compile(exp)
	if err != nil {
		return nil
	}

	for _, res := range *tmpResSlice {
		s := re.FindAllString(res, -1)
		for _, i := range s {
			tmp = append(tmp, i)
		}
	}
	return tmp
}

//get subdomain regular expression
func getExp(field string) string {
	return fmt.Sprintf(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z]{0,62})*\.(%s)$?`, field)
}
