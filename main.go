package main

import (
	"flag"
	"net/http"
)

var pastebin *handler
var conf config
var confpath string

func init() {
	flag.StringVar(&confpath, "c", "./config.json", "Set the config path")

	flag.Parse()

	err := readConf(confpath, &conf)
	if err != nil {
		panic("Open Config Error")
	}

	pastebin = NewHandler(conf.BuffSize)

	go pastebin.timeToCleanUp(conf.CleanDur)
}

func main() {
	http.HandleFunc("/add", pastebin.add)
	http.HandleFunc("/get", pastebin.get)
	http.HandleFunc("/del", pastebin.del)

	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
	}
}
