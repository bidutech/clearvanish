package config

/**************************
author:liyuduo
date:2016.06.01
*********************************/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	//"path"
	"path/filepath"
	"strings"
)

type Host struct {
	Host string `json:"host"` //127.0.0.1
	Port int    `json:"port"` //8000
}

type VanishServers struct {
	Local   Host   `json:"local"` //127.0.0.1:8080
	Servers []Host `json:"servers"`
}

var VanishServer VanishServers

const (
	serversConfigPath = "servers.json"
)

func GetCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path = string(path[0:(strings.LastIndex(path, "/") + 1)])
	return path
}

func InitConfigSrc(configpath string, conf interface{}) {
	f, err := os.Open(configpath)
	defer f.Close()
	if nil == err {
		buff := bufio.NewReader(f)
		for {
			line, err := buff.ReadBytes('\n')
			if err != nil || io.EOF == err {
				return
			}
			errjson := json.Unmarshal(line, &conf)
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

	curentpath := GetCurrentPath()
	serverspath := curentpath + serversConfigPath
	fmt.Println("serverspath", serverspath)
	InitConfigServers(serverspath, &VanishServer)
	//PrintConfigInfo()

}
