package main

import (
	"net/http"
	"time"

	"github.com/Muskchen/toolkits/holmes"
)

func init() {
	http.HandleFunc("/makelgb", makelgbslice)
	go http.ListenAndServe(":10003", nil)
}

func makelgbslice(w http.ResponseWriter, r *http.Request) {
	var a = make([]byte, 1073741824)
	_ = a
}

func main() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("2s"),
		holmes.WithCoolDown("1m"),
		holmes.WithDumpPath("/tmp"),
		holmes.WithTextDump(),
		holmes.WithMemDump(3, 25, 80))
	h.EnableMemDump().Start()
	time.Sleep(1 * time.Hour)
}
