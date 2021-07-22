package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type handler struct {
	lhm       *LinkedHashMap
	lock      sync.Mutex
	whitelist map[string]struct{}
}

type unit struct {
	text string
	exp  int64
}

func NewHandler(length int) *handler {
	return &handler{
		lhm:       NewLHM(length),
		lock:      sync.Mutex{},
		whitelist: make(map[string]struct{}),
	}
}

func (h *handler) add(w http.ResponseWriter, r *http.Request) {
	ok := h.check(w, r)
	if r.Method != "POST" || !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength > int64(conf.MaxLength)*1024+64 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	tmp := r.PostFormValue("exp")
	if tmp == "" {
		tmp = "1440"
	}
	exp, err := strconv.Atoi(tmp)
	text := r.PostFormValue("text")
	if text == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &unit{
		text: text,
		exp:  time.Now().Unix() + int64(exp*60),
	}

	key := hash(*data)
	h.lock.Lock()
	h.lhm.Add(key, data)
	h.lock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "token_" + key,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Value:    key,
		MaxAge:   exp * 60,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	key, ok0 := r.URL.Query()["k"]
	ok1 := h.check(w, r)
	if !ok1 || !ok0 {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.lock.Lock()
	tmp0, ok := h.lhm.Get(key[0])
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
	key, ok0 := r.URL.Query()["k"]
	ok1 := h.check(w, r)
	cookie, err := r.Cookie("token_" + key[0])
	if !ok0 || !ok1 || err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	h.lock.Lock()
	ok0 = h.lhm.Delete(key[0])
	h.lock.Unlock()

	if !ok0 {
		http.SetCookie(w, &http.Cookie{
			Name:    cookie.Name,
			Expires: time.Unix(1, 0),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    cookie.Name,
		Expires: time.Unix(1, 0),
	})
	w.WriteHeader(http.StatusOK)
}

func (h *handler) check(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	_, ok := h.whitelist[origin]
	if !ok {
		return false
	}
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	return true
}

func (h *handler) cleanup() {
	var data *unit
	tmp := h.lhm.head
	for {
		if tmp == nil {
			return
		}
		data = tmp.data.(*unit)
		if data.exp <= time.Now().Unix() {
			h.lock.Lock()
			h.lhm.Delete(tmp.key)
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
