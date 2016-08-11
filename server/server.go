package server

/**************************
author:liyudo
date:2016.06.01
*********************************/

import (
	"clearvanish/config"
	"clearvanish/controller"
	"fmt"
	"log"
	"net"
	"net/http"
	//"sync"
	//"time"

	"github.com/gorilla/mux"
)

func Tcpserver() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	checkError(err)
	fmt.Println("建立成功!")
	for {
		conn, err := listener.Accept()
		checkError(err)
		go doServerStuff(conn)
	}
}

//处理客户端消息
func doServerStuff(conn net.Conn) {
	nameInfo := make([]byte, 512) //生成一个缓存数组
	_, err := conn.Read(nameInfo)
	checkError(err)

	for {
		buf := make([]byte, 512)
		_, err := conn.Read(buf) //读取客户机发的消息
		flag := checkError(err)
		if flag == 0 {
			break
		}
		fmt.Println(string(buf)) //打印出来
	}
}

//检查错误
func checkError(err error) int {
	if err != nil {
		if err.Error() == "EOF" {
			//fmt.Println("用户退出了")
			return 0
		}
		log.Fatal("an error!", err.Error())
		return -1
	}
	return 1
}

func Httpserver() {
	r := mux.NewRouter()
	controller.Handle(r)
	localhost := fmt.Sprintf("%s:%d", config.VanishServer.Local.Host, config.VanishServer.Local.Port)
	fmt.Println("localhost", localhost)
	http.ListenAndServe(localhost, r)
}
