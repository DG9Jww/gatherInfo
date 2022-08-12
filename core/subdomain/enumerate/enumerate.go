/*
CopyRight 2022
Author:DG9J
*/
package enumerate

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"math/rand"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/core/subdomain/result"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket/pcap"
	"golang.org/x/time/rate"
)

//
var (
	timeout  time.Duration = -1 * time.Second
	snapshot int32         = 65535
	promisc  bool          = false
	myHandle *pcap.Handle
	myEthTab ethTable

	//signal for starting sending DNS packets
	sendingSignal = make(chan bool)

	bruteResults = make(chan RecvResults, 100)

	//tables need being removed
	removedTabChan chan TabInfo

	//the number of total valid subdomain
	total int32
)

type bruter struct {
	domain []string

	//interface information
	ethTab ethTable

	//resolvers' IP address
	resolvers []string

	//pcap handle which is for sending packet
	handle *pcap.Handle

	//from 40000 to 60000,self-increasing
	srcPort uint16

	//rate of sending packet,unit is byte
	rate int64

	//a including domain chan and waiting for retry
	retryChan chan *statusTable

	//record packets status
	statusTabLinkList *tableLinkList

	//for wildcard domain name
	blackList []string
}

//interface information
type ethTable struct {
	devName string
	srcMAC  net.HardwareAddr
	srcIP   net.IP
	dstMAC  net.HardwareAddr
}

func newBruter(cfg *config.SubDomainConfig) *bruter {
	dev := common.AutoGetDevice()
	src_mac, _ := net.ParseMAC(dev["srcMAC"])
	dst_mac, _ := net.ParseMAC(dev["dstMAC"])
	myEthTab = ethTable{
		devName: dev["devName"],
		srcMAC:  src_mac,
		srcIP:   net.ParseIP(dev["srcIP"]),
		dstMAC:  dst_mac,
	}

	packetSize := int64(100) //the size of DNS packet is about 74
	myRate := cfg.BandWidth / packetSize
	myLinkList := initTabLinkList()
	removedTabChan = make(chan TabInfo, myRate)

	myHandle, _ = pcap.OpenLive(myEthTab.devName, snapshot, promisc, timeout)
	b := &bruter{
		domain:    cfg.Domain,
		ethTab:    myEthTab,
		resolvers: []string{"223.5.5.5", "223.6.6.6", "180.76.76.76", "119.29.29.29", "114.114.115.115"},
		//resolvers: []string{"8.8.8.8", "1.0.0.1", "8.8.4.4", "1.1.1.1"},
		handle:            myHandle,
		srcPort:           40000,
		rate:              myRate,
		retryChan:         make(chan *statusTable, myRate),
		statusTabLinkList: myLinkList,
	}
	return b
}

func Run(cfg *config.SubDomainConfig) {
	bruter := newBruter(cfg)
	defer bruter.handle.Close()

	//limit the rate according to the bandwith option
	limiter := rate.NewLimiter(rate.Every(time.Duration(time.Second.Nanoseconds()/bruter.rate)), int(bruter.rate))
	ctx := context.Background()

	//load dictionary
	file := common.LoadFile(cfg.BruteDict)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var recvEndSignal = make(chan struct{})

	//     ===== a goroutine for receiving DNS packet =====
	go bruter.recvDNS(sendingSignal, recvEndSignal)

	//     ===== a goroutine for check timeout packet =====
	//Throught continuous cycle detection,put the time out statusTable into retryChan

	go bruter.checkTimeout(recvEndSignal)

	//    	===== A goroutine for retry =====
	//Trying to get statusTable from retryChan and
	//send packet one more time
	go func(l *rate.Limiter) {
		for table := range bruter.retryChan {
			var resolver string
			for {
				resolver = bruter.getResolver()
				if resolver == table.resolver {
					continue
				}
				break
			}
			l.Wait(ctx)
			bruter.sendDNS(table.domain, resolver, table.flagID)
			table.status = 1
			table.retry++
			table.time = time.Now()
		}
		return
	}(limiter)

	//       ============ detect wildcard domain name ===============
	<-sendingSignal
	logger.ConsoleLog(logger.NORMAL, "Detecting WildCard Domain Name......")
	var domainCounter int
	for index, mainDomain := range bruter.domain {
		//determine whether the domain name is wildcard domain
		ok, blackList := bruter.isWildCard(mainDomain)
		if ok {
			if cfg.WildCard {
				logger.ConsoleLog(logger.WARN, fmt.Sprintf("Detected Wildcard Domain: [%s]  BlackList: %v ", mainDomain, blackList))
			} else {
				logger.ConsoleLog(logger.WARN, fmt.Sprintf("Detected Wildcard Domain: [%s] ,Skip!", mainDomain))
				//remove item because wildcard
				bruter.domain = common.DeleteStringFromSlice(bruter.domain, index)
				continue
			}
		}
		domainCounter++
	}
	if domainCounter == 0 {
		return
	}

	//     ======== a goroutine for processing results ========
	go func(filterWildCard bool) {
		for res := range bruteResults {
			var printer string
			for _, record := range res.records {
				if filterWildCard {
					if bruter.checkBlackList(record) {
						continue
					} else {
						printer += " => " + record
					}
				} else {
					printer += " => " + record
				}
			}

			if printer != "" {
				var r = &result.Result{}
				r.SetSubdomain(res.subdomain)
				r.SetRecord(printer)
				result.FinalResults <- r
				atomic.AddInt32(&total, 1)
			}
		}
	}(cfg.WildCard)

	//     ============ a goroutine for removing statusTable ==============
	go func() {
		for tabInfo := range removedTabChan {
			tab, err := bruter.statusTabLinkList.queryStatusTab(tabInfo.subdomain, tabInfo.flagID)
			if err != nil {
				continue
			}
			bruter.statusTabLinkList.remove(tab)
		}
		return
	}()

	//     ============ sending packets ==============
	if len(bruter.domain) == 0 {
		return
	}
	for scanner.Scan() {
		for _, mainDomain := range bruter.domain {
			//get parameters
			domain := scanner.Text() + "." + mainDomain
			resolver := bruter.getResolver()
			flagID := getFlagID()

			//limit rate
			limiter.Wait(ctx)
			//record status and send DNS packet
			table := bruter.recordStatus(domain, resolver, bruter.srcPort, flagID)
			bruter.sendDNS(domain, resolver, flagID)
			table.status = 1
		}
	}

	<-recvEndSignal
	//<-recvDone
	logger.ConsoleLog(logger.CustomizeLog(logger.GREEN, ""), fmt.Sprintf("===== %d Subdomain Found =====", total))

}

//Get Random Resolver
func (bru *bruter) getResolver() string {
	return bru.resolvers[rand.Intn(len(bru.resolvers))]
}

//Get FlagID
func getFlagID() uint16 {
	return uint16(common.RandomInt64(5000, 6000))
}

//record status on statusTable
func (bru *bruter) recordStatus(domain, resolver string, srcPort uint16, flagID uint16) *statusTable {
	tab := &statusTable{domain: domain, retry: 0, time: time.Now(), status: 0, resolver: resolver, flagID: flagID}
	bru.statusTabLinkList.append(tab)
	return tab
}

//check the timeout item from statusTableChan
//and put the timeout item into retryChan
func (bru *bruter) checkTimeout(recvEndSignal chan struct{}) {
	currentTab := bru.statusTabLinkList.head
	for {

		//invalid
		if currentTab == nil {
			currentTab = bru.statusTabLinkList.head
			continue
		}

		if currentTab.retry >= 2 {
			nextTab := currentTab.next
			err := bru.statusTabLinkList.remove(currentTab)
			//if err equal emptyLink which means the task was finished
			if err == emptyLink {
				close(recvEndSignal)
				return
			}
			currentTab = nextTab
			continue
		}

		//retry
		if time.Since(currentTab.time) > time.Second*5 && currentTab.retry < 2 && currentTab.status == 1 {
			currentTab.status = 0
			bru.retryChan <- currentTab
		}
		currentTab = currentTab.next
	}
}
