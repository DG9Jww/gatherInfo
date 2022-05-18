/*
CopyRight 2022
Author:DG9J
*/

package results

import (
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
)

type SubDomainResults struct {
	//results from different modules
	tempDomain []string
	tempIP     []string

	//final result which includes valid item
	domainList []string
	ipList     []string

	//this chan is for dirScan module which includes subdomains
	domainChan chan string

	//this chan is for portscan module
	ipChan chan string

	lock sync.Mutex
}

func (s *SubDomainResults) AddDomain(d string) {
	s.lock.Lock()
	s.tempDomain = append(s.tempDomain, d)
	s.lock.Unlock()
}
func (s *SubDomainResults) AddIP(i string)             { s.tempIP = append(s.tempIP, i) }
func (s *SubDomainResults) GetDomain() []string        { return s.tempDomain }
func (s *SubDomainResults) GetIP() []string            { return s.tempIP }
func (s *SubDomainResults) GetDomainList() []string    { return s.domainList }
func (s *SubDomainResults) GetIPList() []string        { return s.ipList }
func (s *SubDomainResults) GetDomainChan() chan string { return s.domainChan }
func (s *SubDomainResults) GetIPChan() chan string     { return s.ipChan }

//return a new *SubDomainResults
func NewResults() *SubDomainResults {
	return &SubDomainResults{
		domainChan: make(chan string),
		ipChan:     make(chan string),
	}
}

// verify domain task
func (r *SubDomainResults) VerifyDomainTask(domain string, index int, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		//DO REQUEST
		url := "http://" + domain
		req, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		_, err = common.HttpRequest(req)
		if err != nil {
			return
		}

		//Filter out the valid domain
		logger.ConsoleLog(logger.SUBDOMAIN, domain)
		r.lock.Lock()
		r.domainList = append(r.domainList, domain)
		r.lock.Unlock()
	}
}

//mark
//this one is for situation that dirscan module is set
func (r *SubDomainResults) VerifyDomainTask2(domain string, index int, wg *sync.WaitGroup) func() {
	return func() {
		defer wg.Done()
		//DO REQUEST
		url := "http://" + domain
		req, err := common.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		_, err = common.HttpRequest(req)
		if err != nil {
			return
		}

		//Filter out the valid domain
		logger.ConsoleLog(logger.SUBDOMAIN, domain)
		r.domainChan <- domain
	}
}

// verify domain
func (r *SubDomainResults) VerifyDomain(isDir bool) {
	p := common.NewPool(50)
	var wg sync.WaitGroup
	defer p.Release()

	if !isDir {
		for k, v := range r.tempDomain {
			p.Submit(r.VerifyDomainTask(v, k, &wg))
			wg.Add(1)
		}
		wg.Wait()
	} else {
		for k, v := range r.tempDomain {
			p.Submit(r.VerifyDomainTask2(v, k, &wg))
			wg.Add(1)
		}

		//add ip to ipChan
		boo := make(map[string]bool)
		for _, ip := range r.tempIP {
			if _, ok := boo[ip]; !ok {
				boo[ip] = true
				r.ipChan <- ip
			}
		}
		wg.Wait()

		//close chans
		close(r.domainChan)
		close(r.ipChan)
	}
	logger.ConsoleLog(logger.NORMAL, "SUBDOMAIN COMPLETED")
}

//remove duplicate
func (r *SubDomainResults) RemoveDuplicate() {
	r.FixDomain()
	boo := make(map[string]bool)
	boo2 := make(map[string]bool)
	var newDomainList []string
	var newIPList []string
	for _, domain := range r.tempDomain {
		if _, ok := boo[domain]; !ok {
			boo[domain] = true
			newDomainList = append(newDomainList, domain)
		}
	}
	for _, ip := range r.tempIP {
		if _, ok := boo2[ip]; !ok {
			boo[ip] = true
			newIPList = append(newIPList, ip)
		}
	}
	r.tempDomain = newDomainList
	r.tempIP = newIPList
}

//Drop Suffix such as port number
func (r *SubDomainResults) FixDomain() {
	//trim space
	for index, domain := range r.tempDomain {
		r.tempDomain[index] = strings.TrimSpace(domain)
		//if exist suffix
		if strings.Contains(domain, ":") {
			i := strings.Index(domain, ":")
			r.tempDomain[index] = domain[:i]
		}
	}
}
