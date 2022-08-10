package scripts

import (
	"encoding/json"
)

//crt.sh
type Crt struct {
    Subdomain string `json:"name_value"`
}

//return subdomain slice
func (c Crt) SpecialProcess(b []byte) (s []string,err error) {
    var t []Crt

    err = json.Unmarshal(b,&t)
    if err != nil {
        return nil,err
    }

    for _,v := range t {
        s = append(s, v.Subdomain)
    }
    return s,nil
}
