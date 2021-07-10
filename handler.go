package main

import (
	"net/http"
	"strconv"
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
	if r.Method != "POST" || r.ParseForm() != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tmp := r.Form["exp"][0]
	if tmp == "" {
		tmp = "1440"
	}
	exp, err := strconv.Atoi(tmp)
	text, ok := r.Form["text"]
	if !ok || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &unit{
		text: text[0],
		exp:  time.Now().Unix() + int64(exp*60),
	}

	key := hash(*data)
	h.lock.Lock()
	h.data.Add(key, data)
	h.lock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:   "token_" + key,
		Value:  key,
		MaxAge: exp * 60,
	})
	http.Redirect(w, r, conf.Frontend+"/get?k="+key, http.StatusFound)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	key, ok := r.URL.Query()["k"]
	origin := r.Header.Get("Origin")
	if !ok || origin != conf.Frontend {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")

	h.lock.Lock()
	tmp0, ok := h.data.Get(key[0])
	h.lock.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tmp := tmp0.(*unit)

	text := tmp.text
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (h *handler) del(w http.ResponseWriter, r *http.Request) {
	key, ok := r.URL.Query()["k"]
	cookie, err := r.Cookie("token_" + key[0])
	if !ok || err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You Cant Delete This Paste"))
		return
	}

	h.lock.Lock()
	ok = h.data.Delete(key[0])
	h.lock.Unlock()

	if !ok {
		http.SetCookie(w, &http.Cookie{
			Name:   cookie.Name,
			MaxAge: 0,
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   cookie.Name,
		MaxAge: 0,
	})
	http.Redirect(w, r, conf.Frontend, http.StatusFound)
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
	timer := time.NewTimer(time.Second * time.Duration(dur))
	for {
		<-timer.C
		h.cleanup()
		timer.Reset(time.Second * time.Duration(dur))
	}
}
