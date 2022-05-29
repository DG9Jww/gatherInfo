package portscan

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ethTable struct {
	devName string
	srcMAC  net.HardwareAddr
	srcIP   net.IP
	dstMAC  net.HardwareAddr
}

var (
	option = gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
)

//ACK Scan
func (cli *client) sendACK(port int, ip string, handle *pcap.Handle, ethTab *ethTable) {
	dstIP := net.ParseIP(ip)
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
		SrcPort: layers.TCPPort(49000 + rand.Intn(5000)), // Random source port and used to determine recv dst port range
		DstPort: layers.TCPPort(port),
		ACK:     true,
		Window:  14600,
		Seq:     uint32(500000 + rand.Intn(5000)),
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf, option, tcpLayer)

	conn, err := net.Dial("ip4:tcp", ip)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}
	defer conn.Close()
	conn.Write(buf.Bytes())

	//set deadline and we don't need to wait forever
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		fmt.Println(err)
	}

	for {
		b := make([]byte, 4096)
		n, err := conn.Read(b)
		if err != nil {
			return
		}
		if n > 0 {
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.SrcPort == layers.TCPPort(port) {
					if tcp.SYN && tcp.ACK {
						logger.ConsoleLog(logger.PORTSCAN, fmt.Sprintf("%s ===> %d    open", ip, port))
					}
					return
				}
			}
		}
	}
}

//TCP Scan
func (cli *client) sendTCP(port int, ip string) {
	address := fmt.Sprintf("%s:%d", ip, port)
	_, err := net.DialTimeout("tcp", address, time.Second*3)
	if err == nil {
		cli.lock.Lock()
		cli.result[ip] = append(cli.result[ip], port)
		logger.ConsoleLog(logger.PORTSCAN, fmt.Sprintf("%s ===> %d    open", ip, port))
		cli.lock.Unlock()
	}

}

//it's garbage
func (cli *client) sendUDP(port int, ip string) {

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})
	fmt.Println(ip, port)
	if err != nil {
		fmt.Println(err)
	}
	payload := common.GetRandomPayload()
	conn.Write(payload)
}

//receive packets
func (cli *client) recvPackets(devName string) {
	var (
		timeout  time.Duration = -1 * time.Second
		snapshot int32         = 65535
		promisc  bool          = true
		handle   *pcap.Handle
	)

	handle, _ = pcap.OpenLive(devName, snapshot, promisc, timeout)
	defer handle.Close()
	err := handle.SetBPFFilter("tcp and dst portrange 50000-55000")
	if err != nil {
		logger.ConsoleLog(logger.ERROR, fmt.Sprintf("SetBPFFilter Failed:", err.Error()))
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var eth layers.Ethernet
	var ipv4 layers.IPv4
	var tcp layers.TCP

	//layer paser
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ipv4, &tcp)
	//layers packets will be put into this slice when successfully pasing packet
	var decodedLayers []gopacket.LayerType

	scanMode := cli.scanMode

	//SYN scan
	if scanMode == "sS" {
		for {
			//get packet
			packet, err := packetSource.NextPacket()
			if err != nil {
				continue
			}

			//try parse
			err = parser.DecodeLayers(packet.Data(), &decodedLayers)
			if err != nil {
				continue
			}

			//find the packet we need
			if common.IsStringInSlice(ipv4.SrcIP.String(), cli.ipList) && common.MatchInt(int(tcp.SrcPort), cli.portList) {
				//ACK && SYN which means port opening
				if tcp.ACK && tcp.SYN {
					logger.ConsoleLog(logger.PORTSCAN, fmt.Sprintf("%s ===> %d    open", ipv4.SrcIP.String(), tcp.SrcPort))
				}
			}
		}
	}
}
