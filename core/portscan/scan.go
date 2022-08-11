package portscan

import (
	"net"
	"sync"
	"time"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/config"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket/pcap"
)

func Run(cfg *config.PortScanConfig, wg *sync.WaitGroup) {
	if cfg.Enabled {
		if check(cfg) {
			cli := newClient(cfg)
			cli.Run()
		}
	}
	wg.Done()
}

func (cli *client) Run() {
	logger.ConsoleLog(logger.NORMAL, "PortScan is Running......")
	dev := common.AutoGetDevice()
	src_mac, _ := net.ParseMAC(dev["srcMAC"])
	dst_mac, _ := net.ParseMAC(dev["dstMAC"])
	ethTab := &ethTable{
		devName: dev["devName"],
		srcMAC:  src_mac,
		srcIP:   net.ParseIP(dev["srcIP"]),
		dstMAC:  dst_mac,
	}

	//listen on interface and analyse the specific packets
	go cli.recvPackets(ethTab.devName)

	var (
		snapshot int32 = 65535
		promisc  bool  = true
	)
	//set time out to 1 microsecond to fix some bug
	handle, _ := pcap.OpenLive(ethTab.devName, snapshot, promisc, time.Microsecond*1)
	defer handle.Close()
	//	go cli.recvPackets(handle)
	//start task
	pool := common.NewPool(cli.coroutine)
	defer pool.Release()
	for _, ip := range cli.ipList {
		var wg sync.WaitGroup
		for _, port := range cli.portList {
			pool.Submit(cli.task(&wg, port, ip, handle, ethTab))
			wg.Add(1)
		}
		wg.Wait()
	}
	time.Sleep(time.Second * 15)
	logger.ConsoleLog(logger.NORMAL, "PortScan is Completed")
}

func (cli *client) task(wg *sync.WaitGroup, port int, ip string, handle *pcap.Handle, ethTab *ethTable) func() {
	switch cli.scanMode {
	case "sT":
		return func() { cli.sendTCP(port, ip); wg.Done() }
	case "sS":
		return func() { //cli.sendSYN(port, ip, handle, ethTab); wg.Done()
			cli.synScan(port, ip, handle, ethTab)
			wg.Done()
		}
	case "sA":
		return func() { cli.sendACK(port, ip, handle, ethTab); wg.Done() }
	case "sU":
		return func() { cli.sendUDP(port, ip); wg.Done() }
	}
	return nil
}
