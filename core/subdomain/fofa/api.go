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
	for _, domain := range cli.url {
		data := fmt.Sprintf(`domain="%s"`, domain)
		a, err := cli.doRequest(data, 1)
		if err != nil {
			return err
		}

		//Calculating the sum
		sum := a.Size / 100
		if sum > 0 {
			for i := 2; i <= (sum + 1); i++ {
				cli.doRequest(data, i)
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

//request to fofa
func (cli *Client) doRequest(data string, page int) (a APIInfo, err error) {
	resp, err := http.Get(fmt.Sprintf("%s?email=%s&key=%s&qbase64=%s&size=100&page=%d", BaseUrl, cli.email, cli.key, base64.StdEncoding.EncodeToString([]byte(data)), page))
	if err != nil {
		return a, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return a, err
	}

	for _, v := range a.Results {
		temp := RemovePrefix(v[0])
		cli.SubDomainResults.AddDomain(temp)
		cli.SubDomainResults.AddIP(v[1])
	}
	return a, nil
}
