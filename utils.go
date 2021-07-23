package main

import (
	"encoding/hex"
	"hash/fnv"
	"strconv"
	"strings"
)

func hash(data paste) string {
	tmp := []byte(data.text + strconv.Itoa(int(data.exp)))
	hasher := fnv.New64a()
	hasher.Write(tmp)
	key := hex.EncodeToString(hasher.Sum(nil))
	return strings.ToUpper(key)
}
