/*
CopyRight 2022
Author:DG9J
*/
package brute

import (
	"bufio"
	"context"
	"net"
	"time"

	"math/rand"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/google/gopacket/pcap"
	"golang.org/x/time/rate"
)

var (
	timeout  time.Duration = -1 * time.Second
	snapshot int32         = 65535
	promisc  bool          = true
	myHandle *pcap.Handle
)

type bruter struct {
	domain string

	//interface information
	ethTab ethTable

	//resolvers' IP address
	resolvers []string

	//pcap handle which is for sending packet
	handle *pcap.Handle

	//from 40000 to 60000,self-increasing
	srcPort uint16

	//rate of sending packet
	rate int64

<<<<<<< HEAD
	//save statusTable
	//statusTabChan chan statusTable

	//a including domain chan and waiting for retry
	retryChan chan *statusTable

	//
	statusTabList []statusTable
=======
	//
	statusTabChan chan statusTable
>>>>>>> a974f4de79b9a00725fb2cdad7a1d25100f010e0
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

	//sending packet time
	time time.Time
<<<<<<< HEAD
=======

	//last resolver used
	resolver string
>>>>>>> a974f4de79b9a00725fb2cdad7a1d25100f010e0

	//the amount of attempts
	retry int8
}

func newBruter(bandwith int64) *bruter {
	dev := common.AutoGetDevice()
	src_mac, _ := net.ParseMAC(dev["srcMAC"])
	dst_mac, _ := net.ParseMAC(dev["dstMAC"])
	myEthTab := ethTable{
		devName: dev["devName"],
		srcMAC:  src_mac,
		srcIP:   net.ParseIP(dev["srcIP"]),
		dstMAC:  dst_mac,
	}

	packetSize := int64(100) //the size of DNS packet is about 74
	myRate := bandwith / packetSize
	myHandle, _ = pcap.OpenLive(myEthTab.devName, snapshot, promisc, timeout)
	b := &bruter{
		ethTab:    myEthTab,
		resolvers: []string{"223.5.5.5", "223.6.6.6", "180.76.76.76", "119.29.29.29", "114.114.114.115"},
		handle:    myHandle,
		srcPort:   40000,
		rate:      myRate,
	}
	return b
}

func Run(cfg *config.SubDomainConfig) {
	bruter := newBruter(cfg.BandWith)
	defer bruter.handle.Close()

	//limit the rate according to the bandwith option
	limiter := rate.NewLimiter(rate.Every(time.Duration(time.Second.Nanoseconds()/bruter.rate)), int(bruter.rate))
	ctx := context.Background()

	//load dictionary
	file := common.LoadFile(cfg.BruteDict)
	scanner := bufio.NewScanner(file)

	// a goroutine for retry
	go func() {
<<<<<<< HEAD
		for {
			bruter.checkTimeout()
			for table := range bruter.retryChan {
				resolver := bruter.getResolver()
				flagID := getFlagID()
				bruter.sendDNS(table.domain, resolver, flagID)
				table.retry++
			}
		}
=======
		<-bruter.statusTabChan
>>>>>>> a974f4de79b9a00725fb2cdad7a1d25100f010e0
	}()

	//send packet
	for scanner.Scan() {
		//get parameters
		domain := scanner.Text() + "." + bruter.domain
		resolver := bruter.getResolver()
		flagID := getFlagID()

		//limit rate
		limiter.Wait(ctx)
		//record status and send DNS packet
		bruter.recordStatus(domain, resolver)
		bruter.sendDNS(domain, resolver, flagID)

	}
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
func (bru *bruter) recordStatus(domain, resolver string) {
<<<<<<< HEAD
	tab := statusTable{domain: domain, retry: 0, time: time.Now()}
	bru.statusTabList = append(bru.statusTabList, tab)
}

//check the timeout item from statusTableChan
//and channel the timeout item into retryChan
func (bru *bruter) checkTimeout() bool {
	var signal bool = false
	for _, tab := range bru.statusTabList {
		if time.Since(tab.time) > time.Second*5 && tab.retry < 1 {
			bru.retryChan <- &tab
			signal = true
		}
	}
	return signal
=======
	bru.statusTabChan <- statusTable{domain: domain, resolver: resolver, retry: 0, time: time.Now()}
>>>>>>> a974f4de79b9a00725fb2cdad7a1d25100f010e0
}
