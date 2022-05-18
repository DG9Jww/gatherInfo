/*
CopyRight 2022
Author:DG9J
*/

package config

import (
	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Mode           int
	rootCmd        = &cobra.Command{}
	subDomainCmd   = &cobra.Command{}
	dirScanCmd     = &cobra.Command{}
	portScanCmd    = &cobra.Command{}
	vulScanCmd     = &cobra.Command{}
	fingerPrintCmd = &cobra.Command{}
)

type MyConfig struct {
	SubDomain   SubDomainConfig
	DirScan     DirScanConfig
	PortScan    PortScanConfig
	FingerPrint FingerPrintConfig
	VulScan     VulScanConfig
}

//subdomain config
type SubDomainConfig struct {
	Domain    string
	FofaKey   string
	FofaEmail string
	BandWith  int64
	CensysID  string
	CensysKey string
	BruteDict string
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
	SubDomainInit(&cfg.SubDomain)
	DirScanInit(&cfg.DirScan)
	PortScanInit(&cfg.PortScan)
	FingerPrintInit(&cfg.FingerPrint)

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			switch arg {
			case "subdomain":
				subDomainCmd.Execute()
			case "dirscan":
				dirScanCmd.Execute()
			case "portscan":
				portScanCmd.Execute()
			case "fingerprint":
				fingerPrintCmd.Execute()
			case "vulscan":
				vulScanCmd.Execute()
			}
		}
	}
	rootCmd.Execute()
	return cfg
}

//subdomain command flags
func SubDomainInit(cfg *SubDomainConfig) {
	subDomainCmd = &cobra.Command{
		Use:   "subdomain",
		Short: "Collect for SubDomains",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Enabled = true
			logger.ConsoleLog(logger.NORMAL, "Running SubDomain......")
		},
	}
	subDomainCmd.Flags().StringVarP(&cfg.Domain, "domain", "d", "", "Target Main Domain,such as 'google.com'")
	subDomainCmd.Flags().Int64VarP(&cfg.BandWith, "bandwith", "b", 1000000, "BandWith,unit is byte")
	subDomainCmd.Flags().StringVarP(&cfg.BruteDict, "dict", "p", "dict/dns.txt", "Brute Dictionary Path")
	cfg.FofaKey = viper.GetString("subdomain.fofaKey")
	cfg.FofaEmail = viper.GetString("subdomain.fofaEmail")
	cfg.CensysID = viper.GetString("subdomain.censysID")
	cfg.CensysKey = viper.GetString("subdomain.censysKey")
	rootCmd.AddCommand(subDomainCmd)
}

//dirscan command flags
func DirScanInit(cfg *DirScanConfig) {
	dirScanCmd = &cobra.Command{
		Use:   "dirscan",
		Short: "Dir Scan",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Enabled = true
			logger.ConsoleLog(logger.NORMAL, "Running DirScan......")
		},
	}
	dirScanCmd.Flags().StringVarP(&cfg.UrlDic, "urldict", "U", "", "Url Dictionary Path")
	dirScanCmd.Flags().StringVarP(&cfg.PayloadDic, "payloaddict", "p", "dict/dir.txt", "Payload Dictionary Path")
	dirScanCmd.Flags().StringSliceVarP(&cfg.UrlList, "urls", "u", nil, "Url List(split as ',')")
	dirScanCmd.Flags().IntVarP(&cfg.Coroutine, "thread", "t", 30, "Thread")
	dirScanCmd.Flags().IntSliceVarP(&cfg.ValidCode, "codes", "c", []int{200, 301, 302, 303, 304, 307, 403}, "Valid StatusCode")
	dirScanCmd.Flags().StringVarP(&cfg.FilterStr, "filter", "f", "", "Filter String")
	rootCmd.AddCommand(dirScanCmd)
}

func PortScanInit(cfg *PortScanConfig) {
	var temp string
	portScanCmd = &cobra.Command{
		Use:   "portscan",
		Short: "Port Scan",
		Run: func(cmd *cobra.Command, args []string) {
			if temp != "" {
				cfg.PortList = common.PortToList(temp)
			}
			cfg.Enabled = true
			logger.ConsoleLog(logger.NORMAL, "Running PortScan......")
		},
	}
	portScanCmd.Flags().StringSliceVarP(&cfg.IPList, "iplist", "i", nil, "IP List Readied for Scan")
	portScanCmd.Flags().StringVarP(&temp, "portlist", "p", "", "Port List Readied for Scan")
	portScanCmd.Flags().StringVarP(&cfg.IPDict, "ipdict", "I", "", "IP Dictionary Path")
	portScanCmd.Flags().IntVarP(&cfg.Coroutine, "thread", "t", 100, "Port Scan Thread")
	portScanCmd.Flags().StringVarP(&cfg.Mode, "Mode", "m", "sS", "Port Scan Mode")
	rootCmd.AddCommand(portScanCmd)
}

func VulScanInit(cfg *VulScanConfig) {

}

func FingerPrintInit(cfg *FingerPrintConfig) {
	fingerPrintCmd = &cobra.Command{
		Use:   "fingerprint",
		Short: "FingerPrint Detect",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Enabled = true
			logger.ConsoleLog(logger.NORMAL, "Running SubDomain......")
		},
	}
	fingerPrintCmd.Flags().StringVarP(&cfg.FingerP, "fingerP", "f", "dict/cms.json", "FingerPrint Dictionary Path")
	fingerPrintCmd.Flags().IntVarP(&cfg.Thread, "thread", "t", 100, "FingerPrint Detect")
	fingerPrintCmd.Flags().StringSliceVarP(&cfg.UrlList, "urllist", "u", nil, "Target Url")
	rootCmd.AddCommand(fingerPrintCmd)
}
func init() {
	//Load Configuration File
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.ConsoleLog(logger.ERROR, "Load config file error:", err)
	}
}
