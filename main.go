package main

/**************************
author:liyuduo
date:2016.06.01
*********************************/

import (
	"clearvanish/config"
	//"clearvanish/controller"

	"clearvanish/loger"
	"clearvanish/server"

	"github.com/shanhai2015/SHcommon"

	"fmt"
	"sync"
)

/*
Get方式：l
curl -G -d "spid=10002001&epgid=100106" -H'Host:interface5.voole.com' -H'Purge-ID: 1257'  http://127.0.0.1:8000/
nc -lk 8001

IPlist file name is "varnishiplist.txt"
vanishserver curl -G -d "spid=10002001&epgid=100106" -H'Host:interface5.voole.com' -H'Purge-ID: 1257'  http://172.16.10.216:8000/
*/
func init() {
	loger.Loger.Info("Loger ready")
	config.InitConfig()
}
func Test() {
	var url string
	uri := "/filmlist?id=1414&epid=123"
	port := 80
	for _, varniship := range config.VarnishIpList.IpList {
		if port == 80 {
			url = fmt.Sprintf("http://%s%s", varniship, uri)
		} else {
			url = fmt.Sprintf("http://%s:%d%s", varniship, port, uri)
		}
		fmt.Println(url)
	}
}

func main() {
	fisrt, secend, index := SHcommon.StrSplit("http://sadasdasd.com/split", "sa&&")
	fmt.Println("FIRST", fisrt)
	fmt.Println("SECEND", secend)
	fmt.Println("INDEX", index)
	return

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Httpserver()
	}()
	wg.Wait()
	fmt.Println("work over")

}
