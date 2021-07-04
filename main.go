package main

import (
	"net/http"
)

var pastebin *handler
var conf config

func init() {
	err := readConf("./config.json", &conf)
	if err != nil {
		panic("Open Config Error")
	}

	pastebin = NewHandler(conf.BuffSize)

	go pastebin.timeToCleanUp(conf.CleanDur)
}

func main() {
	http.HandleFunc("/add", pastebin.add)
	http.HandleFunc("/get", pastebin.get)

	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
		if conf.EnableIPv6 {
			http.ListenAndServeTLS("[::]:"+conf.Port,
				conf.CertPath, conf.KeyPath, nil)
		}
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
		if conf.EnableIPv6 {
			http.ListenAndServe("[::]:"+conf.Port, nil)
		}
	}
}
