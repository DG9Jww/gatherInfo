/*
CopyRight 2022
Author:DG9J
*/

package cert

import (
	"github.com/DG9Jww/gatherInfo/core/subdomain/results"
)

var (
	BaseUrl = "https://search.censys.io/api"
)

//client
type Client struct {
	*results.SubDomainResults
	domain string
	apiID  string
	apiKey string
}

func NewClient(r *results.SubDomainResults, d string, id string, k string) *Client {
	return &Client{
		SubDomainResults: r,
		domain: d,
		apiID:  id,
		apiKey: k,
	}
}


