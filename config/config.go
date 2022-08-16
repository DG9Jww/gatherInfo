/*
CopyRight 2022
Author:DG9J
*/

package config

import (
	"fmt"
	"os"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//config mode
	Mode int

	rootCmd    = &cobra.Command{}
	vulScanCmd = &cobra.Command{}
)

var subDomainCmd = &cobra.Command{
	Use:   "subdomain",
	Short: "Collecting SubDomains",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var dirScanCmd = &cobra.Command{
	Use:   "dirscan",
	Short: "Directory Scan",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var portScanCmd = &cobra.Command{
	Use:   "portscan",
	Short: "Port Scan",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var fingerPrintCmd = &cobra.Command{
	Use:   "fingerprint",
	Short: "FingerPrint Detect",
	Run:   func(cmd *cobra.Command, args []string) {},
}

type MyConfig struct {
	SubDomain   SubDomainConfig
	DirScan     DirScanConfig
	PortScan    PortScanConfig
	FingerPrint FingerPrintConfig
	VulScan     VulScanConfig
}

//subdomain config
type SubDomainConfig struct {
	Domain    []string
	BandWidth int64
	WildCard  bool
	Validate  bool
	BruteDict string
	Mode      string
	OutPut    string
	Enabled   bool
}

//dirscan configuration
type DirScanConfig struct {
	UrlList     []string
	UrlDic      string
	PayloadList []string
	PayloadDic  string
	Coroutine   int
	ValidCode   []int
	FilterStr   string
	Proxy       string
	Method      string
	Header      string
	OutPut      string
	Enabled     bool
}

//port scan config
type PortScanConfig struct {
	Enabled   bool
	IPList    []string
	IPDict    string
	PortList  []int
	Mode      string
	Coroutine int
}

//vulscan config
type VulScanConfig struct {
	Enabled bool
}

//fingerprint config
type FingerPrintConfig struct {
	Enabled bool
	Thread  int
	UrlList []string
	FingerP string
}

func ConfigFileInit() *MyConfig {
	cfg := new(MyConfig)
	viper.Unmarshal(cfg)
	return cfg
}

//command mode config initialize
func ConfigCommandInit() *MyConfig {
	cfg := new(MyConfig)
	rootCmd.AddCommand(subDomainCmd)
	rootCmd.AddCommand(dirScanCmd)
	rootCmd.AddCommand(portScanCmd)
	rootCmd.AddCommand(fingerPrintCmd)

	for _, arg := range os.Args {
		switch arg {
		case "subdomain":
			SubDomainInit(&cfg.SubDomain)
		case "dirscan":
			DirScanInit(&cfg.DirScan)
		case "portscan":
			PortScanInit(&cfg.PortScan)
		case "fingerprint":
			FingerPrintInit(&cfg.FingerPrint)
		case "vulscan":
			VulScanInit(&cfg.VulScan)
		}
	}

	rootCmd.Execute()
	return cfg
}

//subdomain command flags
func SubDomainInit(cfg *SubDomainConfig) {
	subDomainCmd.Run = func(cmd *cobra.Command, args []string) {
		cfg.Enabled = true
		logger.ConsoleLog(logger.NORMAL, "Running SubDomain Module......")
	}

	subDomainCmd.Flags().StringSliceVarP(&cfg.Domain, "domain", "d", nil, "Target Main Domain,such as 'google.com'")
	subDomainCmd.Flags().Int64VarP(&cfg.BandWidth, "bandwidth", "b", 30000, "BandWidth,unit is byte. 30000 indicates about 300 packets / second")
	subDomainCmd.Flags().StringVarP(&cfg.BruteDict, "dict", "p", "dict/subdomain.txt", "Payload Dictionary Path For Brute")
	subDomainCmd.Flags().StringVarP(&cfg.Mode, "mode", "m", "", "Subdomain moudule mode")
	subDomainCmd.Flags().StringVarP(&cfg.OutPut, "output", "o", "", "OutPut file path")
	subDomainCmd.Flags().BoolVarP(&cfg.WildCard, "wildcard", "w", false, "Scanning wildCard domain name,default is closed")
	subDomainCmd.Flags().BoolVarP(&cfg.Validate, "validate", "v", false, "Validating the subdomains whether they live")
}

//dirscan command flags
func DirScanInit(cfg *DirScanConfig) {
	dirScanCmd.Run = func(cmd *cobra.Command, args []string) {
		cfg.Enabled = true
		logger.ConsoleLog(logger.NORMAL, "Running DirScan Module......")
	}
	dirScanCmd.Flags().StringVarP(&cfg.UrlDic, "urldict", "U", "", "Url Dictionary Path")
	dirScanCmd.Flags().StringVarP(&cfg.PayloadDic, "dictionary", "d", "dict/dir.txt", "Payload Dictionary Path")
	dirScanCmd.Flags().StringSliceVarP(&cfg.UrlList, "urls", "u", nil, "Url List(split as ',')")
	dirScanCmd.Flags().IntVarP(&cfg.Coroutine, "thread", "t", 30, "Thread")
	dirScanCmd.Flags().IntSliceVarP(&cfg.ValidCode, "codes", "c", nil, "All status codes are valid by default,you can customize valid code like this 200,301")
	dirScanCmd.Flags().StringVarP(&cfg.FilterStr, "filter", "f", "", "Filter String")
	dirScanCmd.Flags().StringVarP(&cfg.Proxy, "proxy", "p", "", "Proxy,such as http://127.0.0.1:8888")
	dirScanCmd.Flags().StringVarP(&cfg.Method, "method", "m", "GET", "HTTP Request Method such as GET,HEAD,etc")
	dirScanCmd.Flags().StringVarP(&cfg.Header, "header", "H", "", `HTTP Request Header.For example:"Authorization: sercretxxxxxx"`)
	dirScanCmd.Flags().StringVarP(&cfg.OutPut, "output", "o", "", "OutPut file path")
}

func PortScanInit(cfg *PortScanConfig) {
	var temp string
	portScanCmd.Run = func(cmd *cobra.Command, args []string) {
		if temp != "" {
			cfg.PortList = common.PortToList(temp)
		}
		cfg.Enabled = true
		logger.ConsoleLog(logger.NORMAL, "Running PortScan Module......")
	}

	portScanCmd.Flags().StringSliceVarP(&cfg.IPList, "iplist", "i", nil, "IP List Readied for Scan")
	portScanCmd.Flags().StringVarP(&temp, "portlist", "p", "", "Port List Readied for Scan")
	portScanCmd.Flags().StringVarP(&cfg.IPDict, "ipdict", "I", "", "IP Dictionary Path")
	portScanCmd.Flags().IntVarP(&cfg.Coroutine, "thread", "t", 100, "Port Scan Thread")
	portScanCmd.Flags().StringVarP(&cfg.Mode, "Mode", "m", "sS", "Port Scan Mode")
}

func VulScanInit(cfg *VulScanConfig) {

}

func FingerPrintInit(cfg *FingerPrintConfig) {
	fingerPrintCmd.Run = func(cmd *cobra.Command, args []string) {
		cfg.Enabled = true
		logger.ConsoleLog(logger.NORMAL, "Running FingerPrint Module......")
	}

	fingerPrintCmd.Flags().StringVarP(&cfg.FingerP, "fingerP", "f", "dict/cms.json", "FingerPrint Dictionary Path")
	fingerPrintCmd.Flags().IntVarP(&cfg.Thread, "thread", "t", 100, "FingerPrint Detect")
	fingerPrintCmd.Flags().StringSliceVarP(&cfg.UrlList, "urllist", "u", nil, "Target Url")
}

func init() {
	//Load Configuration File
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.ConsoleLog(logger.ERROR, fmt.Sprintf("Load config file error:%s", err.Error()))
	}
}
