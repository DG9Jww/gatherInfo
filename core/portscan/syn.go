package portscan

import (
	"math/rand"
	"net"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func (cli *client) sendSyn(handle *pcap.Handle, lays ...gopacket.SerializableLayer) {
	buf := cli.bufPool.Get().(gopacket.SerializeBuffer)
	defer func() {
		buf.Clear()
		cli.bufPool.Put(buf)
	}()

	err := gopacket.SerializeLayers(buf, option, lays...)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
		return
	}
	//flush and send packet
	err = handle.WritePacketData(buf.Bytes())
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
		return
	}
}

func (cli *client) synScan(port int, ip string, handle *pcap.Handle, ethTab *ethTable) {
	dstIP := net.ParseIP(ip)
	payload := common.GetRandomPayload()

	ethLayer := &layers.Ethernet{
		SrcMAC:       ethTab.srcMAC,
		DstMAC:       ethTab.dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		SrcIP:    ethTab.srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      128,
		Id:       uint16(50000 + rand.Intn(500)),
		Flags:    layers.IPv4DontFragment,
		Protocol: layers.IPProtocolTCP,
	}
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(50000 + rand.Intn(5000)), // Random source port and used to determine recv dst port range
		DstPort: layers.TCPPort(port),
		SYN:     true,
		Window:  65280,
		Seq:     uint32(500000 + rand.Intn(5000)),
		Options: []layers.TCPOption{
			{
				OptionType:   layers.TCPOptionKindMSS,
				OptionLength: 4,
				OptionData:   []byte{0x05, 0x50}, // 1360
			},
			{
				OptionType: layers.TCPOptionKindNop,
			},
			{
				OptionType:   layers.TCPOptionKindWindowScale,
				OptionLength: 3,
				OptionData:   []byte{0x08},
			},
			{
				OptionType: layers.TCPOptionKindNop,
			},
			{
				OptionType: layers.TCPOptionKindNop,
			},
			{
				OptionType:   layers.TCPOptionKindSACKPermitted,
				OptionLength: 2,
			},
		},
	}

	//compute sum
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, option, ethLayer, ipLayer, tcpLayer, gopacket.Payload(payload))
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
		return
	}
	//flush and send packet
	err = handle.WritePacketData(buf.Bytes())
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err)
		return
	}
	//for i := 0; i < 3; i++ {
	//	cli.sendSyn(handle, ethLayer, ipLayer, tcpLayer, gopacket.Payload(payload))
	//}
}
