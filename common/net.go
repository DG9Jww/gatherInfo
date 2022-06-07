package common

import (
	"crypto/tls"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var proxy = func(*http.Request) (*url.URL, error) {
	return url.Parse("https://127.0.0.1:8080")
}
var globalTransport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, Proxy: proxy}

//redirect is forbidden
func directPolicyFunc(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func NewHttpClient() *http.Client {
	return &http.Client{
		CheckRedirect: directPolicyFunc,
		Timeout:       time.Second * 15,
		Transport:     globalTransport,
	}
}

func HttpRequest(req *http.Request) (*http.Response, error) {
	cli := NewHttpClient()
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", randomAgent())
	return r, nil
}

func ReadHttpBody(r *http.Response) []byte {
	if r != nil {

		defer r.Body.Close()
		content, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ConsoleLog(logger.WARN, err.Error())
		}
		return content
	}
	return nil
}

//random user agent
func randomAgent() string {
	var headers = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (SymbianOS/9.4; Series60/5.0 NokiaN97-1/20.0.019; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) BrowserNG/7.1.18124",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(headers))
	return headers[n]
}

//get network interface
func AutoGetDevice() map[string]string {
	target := "8.8.8.8"
	signal := make(chan map[string]string)
	var device map[string]string

	//get all networking interfaces
	devices, _ := pcap.FindAllDevs()
	var (
		snapshot int32         = 65535
		promisc  bool          = true
		timeout  time.Duration = time.Second * -1
	)

	//catch packet whose destination is 114.114.114.114
	for _, dev := range devices {
		//filter out valid interfaces in case panic happened
		for _, address := range dev.Addresses {
			devIP := address.IP
			if devIP.To4() == nil || devIP.IsLoopback() {
				continue
			}

			//this coroutine is for catching packet
			go func(deviceName string) {
				handle, err := pcap.OpenLive(deviceName, snapshot, promisc, timeout)
				if err != nil {
					logger.ConsoleLog(logger.ERROR, err.Error())
				}
				defer handle.Close()

				packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
				packetChan := packetSource.Packets()
				for packet := range packetChan {
					ipLayer := packet.Layer(layers.LayerTypeIPv4)
					ethLayer := packet.Layer(layers.LayerTypeEthernet)
					if ipLayer != nil {
						ip, _ := ipLayer.(*layers.IPv4)
						if ip.DstIP.String() == target {
							eth, _ := ethLayer.(*layers.Ethernet)
							devHw := map[string]string{"srcMAC": eth.SrcMAC.String(), "devName": deviceName, "srcIP": devIP.String(), "dstMAC": eth.DstMAC.String()}
							signal <- devHw
						}
					}
				}
			}(dev.Name)
		}
	}

	//
	for {
		var ok bool
		select {
		case device = <-signal:
			ok = true
			break
		default:
			net.DialTimeout("tcp", target+":1", time.Second)
			time.Sleep(time.Millisecond * 50)
		}

		if ok {
			break
		}
	}
	return device
}

//random payload for application layer
func GetRandomPayload() []byte {
	var letters = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	p_length := rand.Intn(len(letters))
	p_content := make([]byte, p_length)
	for i := 0; i < p_length; i++ {
		n := rand.Intn(len(letters))
		p_content[i] = letters[n]
	}
	return p_content
}
