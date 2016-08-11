package config

/**************************
author:liyuduo
date:2016.06.01
*********************************/

import (
	"bufio"
	"clearvanish/loger"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/shanhai2015/SHcommon"
)

type Host struct {
	Host string `json:"host"` //127.0.0.1
	Port int    `json:"port"` //8000
}

type VanishServers struct {
	Local         Host   `json:"local"` //127.0.0.1:8080
	VanishPort    int    `json:"varnishport"`
	ResultSendApi string `json:"resultsendapi"`
}

var VanishServer VanishServers

type Ip struct {
	IpList []string
}

var VarnishIpList Ip

func (ip *Ip) IsIp(ipstr string) bool {

	return SHcommon.IsIp(ipstr)
}

const (
	serversConfigPath = "servers.json"
	varnisiplist      = "varnishiplist.txt"
)

func init() {
	VarnishIpList.IpList = []string{}
	InitIpList()
}
func InitIpList() {

	configpath := SHcommon.GetCurrentPath() + varnisiplist
	f, err := os.Open(configpath)
	defer f.Close()
	if nil == err {
		buff := bufio.NewReader(f)
		for {
			line, _, err := buff.ReadLine()
			if err == io.EOF {
				break
			}
			if len(string(line)) == 0 {
				continue
			}
			if VarnishIpList.IsIp(string(line)) {
				VarnishIpList.IpList = append(VarnishIpList.IpList, string(line))
			} else {
				loger.Loger.Error("ERR Ip", string(line))
			}
		}
	} else {
		loger.Loger.Error("Open iplist file Err")
		os.Exit(-1)
	}
}

func InitConfigServers(configpath string, conf *VanishServers) {
	f, err := os.Open(configpath)
	defer f.Close()
	if nil == err {
		buff := bufio.NewReader(f)
		for {
			line, err := buff.ReadBytes('\n')
			if err != nil || io.EOF == err {
				return
			}
			errjson := json.Unmarshal(line, conf)
			if errjson != nil {
				fmt.Printf("ERR InitConfig Configsr.DataPath:%s,line:%s\n", configpath, line)
				os.Exit(-1)
			}
			break
		}
	} else {
		fmt.Printf("read config error-2")
		os.Exit(-1)
	}
}
func PrintConfigInfo() {

	fmt.Println("servers", VanishServer)
	str, err := json.Marshal(VanishServer)
	if err == nil {
		fmt.Println("servers", string(str))
	}
}
func InitConfig() {
	serverspath := SHcommon.GetCurrentPath() + serversConfigPath
	fmt.Println("serverspath", serverspath)
	InitConfigServers(serverspath, &VanishServer)
	//PrintConfigInfo()
}
