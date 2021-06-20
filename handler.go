package main

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net/http"
	"strconv"
	"strings"
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

func hash(data unit) string {
	tmp := []byte(data.text + strconv.Itoa(int(data.exp)))
	hasher := fnv.New64a()
	hasher.Write(tmp)
	key := hex.EncodeToString(hasher.Sum(nil))
	return strings.ToUpper(key)
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
	fmt.Println(key)
}
