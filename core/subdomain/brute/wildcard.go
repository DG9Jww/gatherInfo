package brute

import (
	"net"

	"github.com/google/gopacket/layers"
)

func (bru *bruter) isWildCard(domain string) (bool, []string) {
	prefix := "donotexistdomain"
	prefix2 := "donotexistdomain2"
	var tmpList []string

	domain = prefix + "." + domain
	_, err := net.LookupIP(domain)
	if err == nil {
		//one more time
		domain = prefix2 + "." + domain
		ip, err := net.LookupIP(domain)
		if err == nil {
			//add ip to blackList
			for _, v := range ip {
				bru.blackList = append(bru.blackList, v.String())
				tmpList = append(tmpList, v.String())
			}
			return true, tmpList
		}
		return false, nil
	}
	return false, nil
}

func (bru *bruter) checkBlackList(ip string) bool {
	for _, v := range bru.blackList {
		if ip == v {
			return true
		}
	}
	return false
}

func getIPFromRecord(record layers.DNSResourceRecord) (ip string) {
	return record.String()
}
