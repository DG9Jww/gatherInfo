package brute

import (
	"math/rand"
	"net"
	"sync/atomic"
	"unsafe"

	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
		Protocol: layers.IPProtocolTCP,
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
	udpLayer.SetNetworkLayerForChecksum(ipLayer)

	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, option, ethLayer, ipLayer, udpLayer, dnsLayer)
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
<<<<<<< HEAD

	//srcPort self-increasing
	if bru.srcPort <= 60000 {
		atomic.AddUintptr((*uintptr)(unsafe.Pointer(&bru.srcPort)), uintptr(1))
	} else {
		bru.srcPort = 40000
	}
=======
	//srcPort self-increasing
	atomic.AddUintptr((*uintptr)(unsafe.Pointer(&bru.srcPort)), uintptr(1))
>>>>>>> a974f4de79b9a00725fb2cdad7a1d25100f010e0
}

//receive dns packets
func (bru *bruter) recvDNS(domain string, resolverIP string, srcPort uint16, flagID uint16) {
}
