package scripts

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DG9Jww/gatherInfo/common"
)

//crt
type crt struct {
    Subdomain string `json:"name_value"` 
}

//return subdomain slice
func StartCrt(domain string) (error,[]string) {
    url := fmt.Sprintf("https://crt.sh/?output=json&q=%s",domain)
    req,err := common.NewRequest("GET",url,nil)
    if err!= nil {
        return err,nil
    }

    resp,err := common.HttpRequest(req)
    if err!= nil {
        return err,nil
    }
    defer resp.Body.Close()
    b,_ := io.ReadAll(resp.Body)

    var tmpSlice []crt
    err = json.Unmarshal(b,&tmpSlice)
    if err!= nil {
        return err,nil
    }
    
    var s []string
    for  _,v := range tmpSlice {
        s = append(s, v.Subdomain) 
    }
    exp := common.GetExp(domain)
    s = common.ProRegularExp(&s,exp)  
    return nil,s 

}
