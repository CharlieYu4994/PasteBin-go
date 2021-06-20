package main

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"sync"
	"time"
)

type handler struct {
	data *LinkedHashMap
	lock sync.Mutex
}

type unit struct {
	text string
	exp  int64
}

func hash(data []byte) string {
	hasher := fnv.New128()
	hasher.Write(data)
	return string(hasher.Sum(nil))
}

func NewHandler(length int) *handler {
	return &handler{
		data: NewLHM(length),
		lock: sync.Mutex{},
	}
}

func (h *handler) add(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" || r.ParseForm() != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	text := r.Form["text"][0]

	data := &unit{
		text: text,
		exp:  time.Now().Unix(),
	}

	key := hash([]byte(text))

	h.data.Add(key, data)
	fmt.Println(key)
}
