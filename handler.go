package main

import (
	"fmt"
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

func NewHandler(length int) *handler {
	return &handler{
		data: NewLHM(length),
		lock: sync.Mutex{},
	}
}

func (h *handler) add(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if r.Method != "POST" || r.ParseForm() != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	text, ok := r.Form["text"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &unit{
		text: text[0],
		exp:  time.Now().Unix() + (24 * 3600),
	}

	key := hash(*data)

	h.data.Add(key, data)
	http.Redirect(w, r, "/get?k="+key, http.StatusFound)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	key, ok := r.URL.Query()["k"]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmp0, ok := h.data.Get(key[0])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tmp := tmp0.(*unit)

	text := tmp.text
	fmt.Println(text)
}
