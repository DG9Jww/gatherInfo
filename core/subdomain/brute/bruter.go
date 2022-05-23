/*
CopyRight 2022
Author:DG9J
*/
package brute

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"math/rand"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket/pcap"
	"golang.org/x/time/rate"
)

var (
	timeout  time.Duration = -1 * time.Second
	snapshot int32         = 65535
	promisc  bool          = false
	myHandle *pcap.Handle
	myEthTab ethTable

	//signal for starting sending DNS packets
	signal = make(chan bool)
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
	statusTabList []*statusTable

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

//status table
type statusTable struct {
	domain string

	srcPort uint16

	//sending packet time
	time time.Time

	//last resolver used
	resolver string

	//the amount of attempts
	retry int8

	//status,0 unsent,1 sent
	status uint8
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

	myHandle, _ = pcap.OpenLive(myEthTab.devName, snapshot, promisc, timeout)
	b := &bruter{
		domain:    cfg.Domain,
		ethTab:    myEthTab,
		resolvers: []string{"223.5.5.5", "223.6.6.6", "180.76.76.76", "119.29.29.29", "114.114.115.115"},
		//resolvers: []string{"8.8.8.8", "1.0.0.1", "8.8.4.4", "1.1.1.1"},
		handle:    myHandle,
		srcPort:   40000,
		rate:      myRate,
		retryChan: make(chan *statusTable, myRate),
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

	//     ===== a goroutine for receive =====
	go func(chan bool) {
		bruter.recvDNS(signal, cfg.WildCard)
	}(signal)

	//     ===== a goroutine for check timeout packet =====
	//Throught continuous cycle detection,put the time out statusTable into retryChan
	go func(chan bool) {
		for {
			bruter.checkTimeout()
		}
	}(signal)

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
			flagID := getFlagID()
			l.Wait(ctx)
			table.srcPort = bruter.srcPort
			bruter.sendDNS(table.domain, resolver, flagID)
			table.status = 1
			table.retry++
			table.time = time.Now()
		}
	}(limiter)

	//     ======== send packet ========
	//detect wildcard domain name
	for _, mainDomain := range bruter.domain {
		ok, blackList := bruter.isWildCard(mainDomain)
		if ok {
			if cfg.WildCard {
				logger.ConsoleLog2(logger.CustomizeLog(logger.YELLOW, "WARNING"), fmt.Sprintf("Detected Wildcard Domain: [%s]  BlackList: %v ", mainDomain, blackList))
			} else {
				logger.ConsoleLog2(logger.CustomizeLog(logger.YELLOW, "WARNING"), fmt.Sprintf("Detected Wildcard Domain: [%s]  BlackList: %v ,Skip!", mainDomain, blackList))
				continue
			}
		}
	}

	//
	<-signal
	for scanner.Scan() {
		for _, mainDomain := range bruter.domain {
			//get parameters
			domain := scanner.Text() + "." + mainDomain
			resolver := bruter.getResolver()
			flagID := getFlagID()

			//limit rate
			limiter.Wait(ctx)
			//record status and send DNS packet
			table := bruter.recordStatus(domain, resolver, bruter.srcPort)
			bruter.sendDNS(domain, resolver, flagID)
			table.status = 1
		}
	}

	time.Sleep(time.Second * 15)
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
func (bru *bruter) recordStatus(domain, resolver string, srcPort uint16) *statusTable {
	tab := statusTable{domain: domain, retry: 0, time: time.Now(), status: 0, resolver: resolver, srcPort: srcPort}
	bru.statusTabList = append(bru.statusTabList, &tab)
	return bru.statusTabList[len(bru.statusTabList)-1]
}

//check the timeout item from statusTableChan
//and channel the timeout item into retryChan
func (bru *bruter) checkTimeout() {
	if bru.statusTabList == nil {
		return
	}
	for _, tab := range bru.statusTabList {
		if time.Since(tab.time) > time.Second*5 && tab.retry < 2 && tab.status == 1 {
			tab.status = 0
			bru.retryChan <- tab
		}
	}
}
