package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	conf "github.com/labs/tracing/config"
)

var (
	method = flag.String("method", "format", "input your method")
)

func main() {
	flag.Parse()
	port, ok := conf.RoleToPort[*method]
	if !ok {
		panic("invalid method")
	}
	log.Printf("start service %s on port %d", *method, port)

	http.HandleFunc("/format", formatHandler)
	http.HandleFunc("/publish", publishHandler)

	url := ":" + strconv.FormatInt(int64(port), 10)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Println("server exit")
	}
}

func formatHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	helloTo := r.FormValue("helloTo")
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)

	w.Write([]byte(helloStr))
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	helloStr := r.FormValue("helloStr")
	log.Println(helloStr)

	w.Write([]byte{})
}
