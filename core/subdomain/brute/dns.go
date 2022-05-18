package brute

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

	fmt.Println(domain)
	//compute sum
	err := udpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}

	buf := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buf, option, ethLayer, ipLayer, udpLayer, dnsLayer)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
		return
	}
	//flush and send packet
	err = bru.handle.WritePacketData(buf.Bytes())
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
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
func (bru *bruter) recvDNS(domain string, resolverIP string, srcPort uint16, flagID uint16) {
	handle, _ := pcap.OpenLive(myEthTab.devName, snapshot, promisc, timeout)
	err := handle.SetBPFFilter("udp and src port 53")
	if err != nil {
		logger.ConsoleLog(logger.ERROR, "SetBPFFilter Failed,", err.Error())
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

	for {
		//get packet
		packet, err := packetSource.NextPacket()
		if err != nil {
			continue
		}

		//trying to parse packet into DNS packet
		err = parser.DecodeLayers(packet.Data(), &decodedLayers)
		if err != nil {
			continue
		}

		//QR = 0 witch means it's a query packet
		//QR = 1 witch means it's a response packet
		if !dnsLayer.QR {
			continue
		}

		//Throught verifying packet type and IP address,
		//we can confirm that it is the packet what we need
		if !common.IsStringInSlice(ipv4Layer.SrcIP.String(), bru.resolvers) {
			continue
		}

		//now start analysing the packet
	}
}
