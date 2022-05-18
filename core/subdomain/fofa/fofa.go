/*
CopyRight 2022
Author:DG9J
*/
package fofa

import (
	"github.com/DG9Jww/gatherInfo/core/subdomain/results"
)

const BaseUrl = "https://fofa.info/api/v1/search/all"

type Client struct {
	*results.SubDomainResults
	url   string
	key   string
	email string
}

func NewClient(r *results.SubDomainResults,k string, e string, u string) *Client {
	return &Client{SubDomainResults: r,key: k, email: e, url: u}
}
