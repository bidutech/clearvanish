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
	"clearvanish/loger"
	//"net/url"
	"strings"

	"sync"

	"github.com/go-httpclient"
	"github.com/gorilla/mux"
	"github.com/shanhai2015/SHcommon"
)

const checkinterval = 5 * 60 //5 minute interval to check

var requeststores map[string]int64
var mutx sync.RWMutex

func init() {
	requeststores = make(map[string]int64)
	go checkResult()
}

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
	loger.Loger.Info("Redirect:", url)
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
		//if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		bodystr = strings.Replace(bodystr, "\n", " ", -1)
		loger.Loger.Info(fmt.Sprintf("Reponse: Url:%s,StatusCode:%d,ResponseBoby:%s", url, response.StatusCode, bodystr))
		//}
	} else {
		loger.Loger.Error(fmt.Sprintf("Redirect: Url:%s,Error:%s", url, err.Error()))
	}
}

func realClear(r *http.Request) {
	var uri string
	uri = r.RequestURI
	n := strings.LastIndex(uri, "/")
	if n > 0 {
		uri = SHcommon.Substr(uri, n, len(uri))
	}
	host := r.Host
	go recordRequest("id")
	port := config.VanishServer.VanishPort
	var url string
	for _, varniship := range config.VarnishIpList.IpList {
		if port == 80 || port == 0 {
			url = fmt.Sprintf("http://%s%s", varniship, uri)
		} else {
			url = fmt.Sprintf("http://%s:%d%s", varniship, port, uri)
		}
		go request(host, url, r)
	}

}

func PrintHeader(r *http.Request) {
	url := r.URL.String()
	host := r.Host
	method := r.Method
	var headerstr string
	for k, v := range r.Header {
		for _, vv := range v {
			headerstr += fmt.Sprintf("%s:%s,", k, vv)
		}
	}
	headerstr = SHcommon.Substr(headerstr, 0, len(headerstr)-1)
	loger.Loger.Info("Request:", fmt.Sprintf("Host:%s,Method:%s,Url:%s,Header:%s", host, method, url, headerstr))
}

func TestReturnResult() {
	returnResult("")
}
func returnResult(requestId string) {
	transport := &httpclient.Transport{
		ConnectTimeout:        10 * time.Second,
		RequestTimeout:        10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", config.VanishServer.ResultSendApi, nil)
	resp, err := client.Do(req)
	if err != nil {
		loger.Loger.Error(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			//do some thing
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			bodystr := string(body)
			loger.Loger.Info(bodystr)
		}
	}

}

func respose() {
	mutx.Lock()
	defer mutx.Unlock()
	now := time.Now().Unix()
	if len(requeststores) > 0 {
		for k, v := range requeststores {
			if (now - v) >= checkinterval {
				go returnResult(k)
			}
		}
	}
}

func checkResult() {
	timer := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer.C:
			respose()
		}
	}
}

func recordRequest(rid string) {
	mutx.Lock()
	defer mutx.Unlock()
	var questId string
	requeststores[questId] = time.Now().Unix()

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
	r.HandleFunc("/{path:.*}", vanaishclear)
	//r.HandleFunc("/", vanaishclear)
}
