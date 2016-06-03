package main

/**************************
author:liyuduo
date:2016.06.01
*********************************/

import (
	"clearvanish/config"
	"clearvanish/loger"
	"clearvanish/server"
	"fmt"
	"sync"
)

func main() {
	loger.InitLog()
	loger.Loger.Info("Loger ready")
	var wg sync.WaitGroup
	config.InitConfig()
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Httpserver()
	}()
	//go server.Tcpserver()
	wg.Wait()
	fmt.Println("work over")
}
