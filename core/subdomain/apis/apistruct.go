package apis

import ()

var APIStruct = map[string]interface{}{
	"virustotal": &virustotal{},
	"censys": &censys{},
}

//virustotal
type virustotal struct {
	Data []virustotal2 `json:"data"`
}

type virustotal2 struct {
	ID string `json:"id"`
}


//censys
type censys struct {
	Results  []censys2 `json:"results"`
}

type censys2 struct {
	Parsed string `json:"parsed.subject_dn"`
}
