package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	EnableTLS bool   `json:"enabletls"`
	BuffSize  int    `json:"buffsize"`
	CertPath  string `json:"certpath"`
	KeyPath   string `json:"keypath"`
	Port      string `json:"port"`
	Frontend  string `json:"frontend"`
}

func readConf(path string, conf *config) error {
	_, err := os.Stat(path)
	if err != nil || os.IsExist(err) {
		return err
	}
	tmp, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, conf)
}
