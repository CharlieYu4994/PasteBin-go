package main

import (
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
		exp:  time.Now().Unix() + (5 * 60),
	}

	key := hash(*data)

	h.data.Add(key, data)
	http.Redirect(w, r, conf.Frontend+"/get?k="+key, http.StatusFound)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	key, ok := r.URL.Query()["k"]
	origin := r.Header.Get("Origin")
	if !ok || origin != conf.Frontend {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")

	tmp0, ok := h.data.Get(key[0])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tmp := tmp0.(*unit)

	text := tmp.text
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (h *handler) cleanup() {
	var data *unit
	tmp := h.data.head
	for {
		if tmp == nil {
			return
		}
		data = tmp.data.(*unit)
		if data.exp <= time.Now().Unix() {
			h.lock.Lock()
			h.data.Delete(tmp.key)
			h.lock.Unlock()
		}
		tmp = tmp.next
	}
}

func (h *handler) timeToCleanUp(dur int) {
	timer := time.NewTimer(0)
	for {
		<-timer.C
		h.cleanup()
		timer.Reset(time.Second * time.Duration(dur))
	}
}
