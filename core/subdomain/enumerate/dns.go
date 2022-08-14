package enumerate

import (
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"unsafe"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	option = gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
)

type RecvResults struct {
	subdomain string
	records   []string
}


//send DNS packet
func (bru *bruter) sendDNS(domain string, resolverIP string, flagID uint16) {
	dstIP := net.ParseIP(resolverIP)
	ethLayer := &layers.Ethernet{
		SrcMAC:       bru.ethTab.srcMAC,
		DstMAC:       bru.ethTab.dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		SrcIP:    bru.ethTab.srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      128,
		Id:       uint16(50000 + rand.Intn(500)),
		Flags:    layers.IPv4DontFragment,
		Protocol: layers.IPProtocolUDP,
	}

	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(bru.srcPort),
		DstPort: layers.UDPPort(53),
	}

	dnsLayer := &layers.DNS{
		ID:      flagID,
		QDCount: 1,
		RD:      false, //recursive query flag
	}

	dnsLayer.Questions = append(dnsLayer.Questions,
		layers.DNSQuestion{
			Name:  []byte(domain),
			Type:  layers.DNSTypeA,
			Class: layers.DNSClassIN,
		})

	//compute sum
	err := udpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}

	buf := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buf, option, ethLayer, ipLayer, udpLayer, dnsLayer)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}
	//flush and send packet
	err = bru.handle.WritePacketData(buf.Bytes())
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
		return
	}

	//srcPort self-increasing
	if bru.srcPort <= 60000 {
		atomic.AddUintptr((*uintptr)(unsafe.Pointer(&bru.srcPort)), uintptr(1))
	} else {
		bru.srcPort = 40000
	}
}

//receive dns packets
func (bru *bruter) recvDNS(signal chan bool, recvEndSignal chan struct{}) {
	handle, _ := pcap.OpenLive(myEthTab.devName, snapshot, promisc, pcap.BlockForever)
	defer handle.Close()
	err := handle.SetBPFFilter("udp and src port 53")
	if err != nil {
		logger.ConsoleLog(logger.ERROR, fmt.Sprintf("SetBPFFilter Failed:%s", err.Error()))
		return
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var ethLayer layers.Ethernet
	var ipv4Layer layers.IPv4
	var ipv6Layer layers.IPv6
	var udpLayer layers.UDP
	var dnsLayer layers.DNS

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet, &ethLayer, &ipv4Layer, &ipv6Layer, &udpLayer, &dnsLayer)
	var decodedLayers []gopacket.LayerType

	//recv goroutine is ready
	close(signal)
	for {
		select {
		case <-recvEndSignal:
			close(bruteResults)
			close(bru.retryChan)
			return
		default:

			/*
				5.19
				Here I got some issues.
				When the network environment is pure,the program was always stuck here
				because here is no more packet and this function can not read any packet.

				How to solve this issues?
				I set the parameter "timeout" to 1 millisecond and use NextPacket to read packet.
				"packetChan := packetSource.Packet()" is doesn't work

				10 mins later......
				Well,I found I didn't get the point.
				The point is that the goroutine for receiving DNS packet is inefficient.
				When my machine receives all DNS response packets,the recvDNS goroutine is only starting.
				So I try to make a trigger before sending packet.When recv goroutine is ready,sending
				packets goroutine start.
			*/
			packet, err := packetSource.NextPacket()
			if err != nil {
				continue
			}

			//trying to parse packet into DNS packet
			err = parser.DecodeLayers(packet.Data(), &decodedLayers)
			if err != nil {
				continue
			}

			//QR = 0 witch means it's a query packet QR = 1 witch means it's a response packet
			if !dnsLayer.QR {
				continue
			}

			//Throught verifying packet type and IP address, we can confirm that it is the packet what we need
			if !common.IsStringInSlice(ipv4Layer.SrcIP.String(), bru.resolvers) {
				continue
			}

			//Reply code equal 0 means no error. Besides,most of time,it will be set to 3 which means the subdomain doesn't exist.
			if dnsLayer.ResponseCode != 0 {
				continue
			}

			//no answer
			if dnsLayer.ANCount == 0 && dnsLayer.ARCount == 0 && dnsLayer.NSCount == 0 {
				continue
			}

			//invalid
			if len(dnsLayer.Questions) == 0 {
				continue
			}

			subdomain := string(dnsLayer.Questions[0].Name)
			if !common.IsSliceWithinStr(subdomain, bru.domain) {
				continue
			}

			//answers
			if dnsLayer.ANCount > 0 {
				//query packet and delete packet
				tmpResult := RecvResults{subdomain: subdomain}
				for _, record := range dnsLayer.Answers {
					tmpResult.records = append(tmpResult.records, record.String())
				}
				bruteResults <- tmpResult

                //query node
				tab, err := bru.statusTabLinkList.queryStatusTab(subdomain, dnsLayer.ID)
				if err != nil {
					continue
				}
                bru.statusTabLinkList.remove(tab)
			}

		}
	}
}
