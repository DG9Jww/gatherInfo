package scripts

import (
	"encoding/json"
)

//fofa
type Fofa struct {
    Results [][]string `json:"results"`
}

//return subdomain slice
func (f Fofa) SpecialProcess(b []byte) (s []string,err error) {
    err = json.Unmarshal(b,&f)
    if err != nil {
        return nil,err
    }

    for _,v := range f.Results {
       s = append(s, v[0])  
    }

    return s,nil
}
