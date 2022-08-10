package apis

import (
	"github.com/DG9Jww/gatherInfo/core/subdomain/apis/scripts"
)

var APIStruct = map[string]interface{}{
	"virustotal":  &virustotal{},
	"censys":      &censys{},
	"threatminer": &threatminer{},
}

type SpecialResp interface {
	SpecialProcess([]byte) ([]string,error)
}

var SpecialRespMap = map[string]SpecialResp{
	"fofa": scripts.Fofa{},
	"crt": scripts.Crt{},
}

//Five field names have been defined,
//you should use them in the API structs so that
//you can parse API successfully.
// 1.[Subdomain]  which means the tag is subdomain results
// 2.[IPaddress]  which means the tag is ipaddress results
// 3.[SubdomainSlice]  which means the tag is subdomain slice results
// 4.[IPaddressSlice]  which means the tag is ipaddress slice results
// 5.[DomainAndIP]  which means ipaddress and subdomain appear in one field
//
//virustotal
type virustotal struct {
	Data []virustotal2 `json:"data"`
}

type virustotal2 struct {
	Subdomain string `json:"id"`
}

//censys
type censys struct {
	Results []censys2 `json:"results"`
}

type censys2 struct {
	Subdomain string `json:"parsed.subject_dn"`
}

//threatminer
type threatminer struct {
	SubdomainSlice []string `json:"results"`
}

