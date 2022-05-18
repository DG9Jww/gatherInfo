/*
CopyRight 2022
Author:DG9J
*/

package fofa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type APIInfo struct {
	Error   bool       `json:"error"`
	Results [][]string `json:"results"`
	Size    int        `json:"size"`
}

func (cli *Client) GetAPIInfo() error {
	data := fmt.Sprintf(`domain="%s"`, cli.url)
	resp, err := http.Get(fmt.Sprintf("%s?email=%s&key=%s&qbase64=%s&size=100&page=1", BaseUrl, cli.email, cli.key, base64.StdEncoding.EncodeToString([]byte(data))))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var a APIInfo
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return err
	}

	for _, v := range a.Results {
		temp := RemovePrefix(v[0])
		cli.SubDomainResults.AddDomain(temp)
		cli.SubDomainResults.AddIP(v[1])
	}

	//Calculating the sum
	sum := a.Size / 100
	//sum > 0 which means result page > 1
	if sum > 0 {
		for i := 2; i <= (sum + 1); i++ {
			resp, err := http.Get(fmt.Sprintf("%s?email=%s&key=%s&qbase64=%s&size=100&page=%d", BaseUrl, cli.email, cli.key, base64.StdEncoding.EncodeToString([]byte(data)), i))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var a APIInfo
			if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
				return err
			}
			for _, v := range a.Results {
				temp := RemovePrefix(v[0])
				cli.SubDomainResults.AddDomain(temp)
				cli.SubDomainResults.AddIP(v[1])
			}
		}
	}
	return nil
}

//remove "http" prefix
func RemovePrefix(a string) string {
	if strings.HasPrefix(a, "http://") {
		b := strings.TrimPrefix(a, "http://")
		return b
	}
	if strings.HasPrefix(a, "https://") {
		b := strings.TrimPrefix(a, "https://")
		return b
	}
	return a
}
