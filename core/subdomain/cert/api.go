/*
CopyRight 2022
Author:DG9J
*/

package cert

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/DG9Jww/gatherInfo/logger"
)

var ErrorCode = map[int]string{
	400: "query could not be parsed",
	404: "page not found",
	429: "rate limit exceeded",
	500: "unknown error occurred",
}

//API Information
type APIResp struct {
	Status   string    `json:"status"`
	Metadata MetaData  `json:"metadata"`
	Results  []Results `json:"results"`
}

type MetaData struct {
	Pages int `json:"pages"`
}

type Results struct {
	Parsed string `json:"parsed.subject_dn"`
}

func (cli *Client) Run() error {
	for _, domain := range cli.domains {
		payload := make(map[string]interface{})
		payload["query"] = domain
		var res_list []string
		payload["page"] = 1
		payload["fields"] = []string{"parsed.subject_dn"}
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", BaseUrl+"/v1/search/certificates", bytes.NewReader(body))
		if err != nil {
			return err
		}

		//According to the document
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cli.apiID, cli.apiKey))))
		c := http.Client{
			Timeout: 15 * time.Second,
		}
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			val, ok := ErrorCode[resp.StatusCode]
			if !ok {
				logger.ConsoleLog(logger.ERROR, "cert moudule ERROR:unknown error,please check your key")
			} else {
				logger.ConsoleLog(logger.ERROR, fmt.Sprintf("cert moudule ERROR:%s", val))
			}
			return nil
		}
		//Parse Response
		var info = new(APIResp)
		err = json.NewDecoder(resp.Body).Decode(info)
		if err != nil {
			return err
		}
		for _, v := range info.Results {
			res := regexp.MustCompile("CN=(.*)").FindStringSubmatch(v.Parsed)[1]
			if strings.Contains(res, "*") {
				continue
			}
			res_list = append(res_list, res)
			cli.SubDomainResults.AddDomain(res)
		}
	}
	return nil
}
