package apis

import ()


var APIStruct = map[string]interface{}{
	"virustotal": &virustotal{},
	"censys": &censys{},
}

//Four field names have been defined,
//you should use them in the API structs so that
//you can parse API successfully.
// 1.[Subdomain]  which means the tag is subdomain results    
// 2.[IPaddress]  which means the tag is ipaddress results    
// 3.[SubdomainSlice]  which means the tag is subdomain slice results    
// 1.[IPaddressSlice]  which means the tag is ipaddress slice results    
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
	Results  []censys2 `json:"results"`
}

type censys2 struct {
	Subdomain string `json:"parsed.subject_dn"`
}
