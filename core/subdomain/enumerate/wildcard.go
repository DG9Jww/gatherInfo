package enumerate

import (
	"time"
)

func (bru *bruter) isWildCard(domain string) (bool, []string) {
	prefix := "donotexistdomain"
	prefix2 := "donotexistdomain2"
	var tmpList []string

	domainName := prefix + "." + domain
	resolver := bru.getResolver()
	flagID := getFlagID()
	bru.sendDNS(domainName, resolver, flagID)
	recvTimeout := time.After(time.Second * 3)
	select {
	case res := <-bruteResults:
		for _, r := range res.records {
			bru.blackList = append(bru.blackList, r)
			tmpList = append(tmpList, r)
		}
		return true, tmpList
	case <-recvTimeout:
		//try again
		domainName := prefix2 + "." + domain
		resolver := bru.getResolver()
		flagID := getFlagID()
		bru.sendDNS(domainName, resolver, flagID)
		recvTimeout := time.After(time.Second * 3)
		select {
		case <-recvTimeout:
			return false, nil
		case res := <-bruteResults:
			for _, r := range res.records {
				bru.blackList = append(bru.blackList, r)
				tmpList = append(tmpList, r)
			}
			return true, tmpList
		}
	}
}

func (bru *bruter) checkBlackList(record string) bool {
	for _, v := range bru.blackList {
		if v == record {
			return true
		}
	}
	return false
}
