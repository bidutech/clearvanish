package controller

/**************************
author:liyuduo
date:2016.06.01
*********************************/

import (
	"clearvanish/config"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	//"strconv"

	"github.com/gorilla/mux"
)

func Test_realClear(header *http.Header) {

	/*
				   for k, v := range r.Header {
		        for _, vv := range v {
		            w.Header().Add(k, vv)
		        }
		    }

			for key, value := range r.Header {
				fmt.Printf("%s->%-10s", key, value)
				reqest.Header.Add(key, value)
			}

	*/

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", "http://www.baidu.com", nil)

	reqest.Header = *header
	//reqest.Header.Set("Connection", "keep-alive")
	fmt.Printf("%s", reqest.Header)
	response, err := client.Do(reqest)
	if err == nil {
		if response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			bodystr := string(body)
			fmt.Println(bodystr)
		}
	}

}
func request(host, url string, r *http.Request) {

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(30 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*10) //设置连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline) //设置读取超时
				return c, nil
			},
		},
	}
	//client := &http.Client{}
	reqest, _ := http.NewRequest("PURGE", url, nil)
	reqest.Header = r.Header
	reqest.Host = host
	//reqest.Method = r.Method
	response, err := client.Do(reqest)
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {

			body, _ := ioutil.ReadAll(response.Body)
			bodystr := string(body)
			fmt.Println(bodystr)
		}
	} else {
		fmt.Print("Time out\n")
	}
}
func realClear(r *http.Request) {
	var uri string
	uri = r.RequestURI
	fmt.Println("uri", uri)
	host := r.Host
	for _, server := range config.VanishServer.Servers {
		//serverhost := fmt.Sprintf("%s:%d", server.Host, server.Port)
		url := fmt.Sprintf("http://%s:%d%s", server.Host, server.Port, uri)
		go request(host, url, r)
	}

}

func PrintHeader(r *http.Request) {
	fmt.Println("http.Request.Host---->", r.Host)
	fmt.Println("http.Request.Method---->", r.Method)
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Println("http.Request.Header:", k, vv)
		}
	}
}

func vanaishclear(w http.ResponseWriter, r *http.Request) {
	PrintHeader(r)
	go realClear(r)
	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(200)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
			<vanishclear>
				<return>200</return>
			</vanishclear>`
	w.Write([]byte(xml))
}

func Handle(r *mux.Router) {
	r.HandleFunc("/", vanaishclear)
}
