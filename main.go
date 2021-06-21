package main

import "net/http"

var pastebin *handler

func init() {
	pastebin = NewHandler(100)
}

func main() {
	http.HandleFunc("/add", pastebin.add)
	http.HandleFunc("/get", pastebin.get)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
