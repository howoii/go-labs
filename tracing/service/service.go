package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labs/tracing/config"
)

var (
	role = flag.String("role", "formatter", "input your role")

	roleToPort = map[string]int32{
		"formatter": 8080,
		"publisher": 8081,
	}
)

func main() {
	flag.Parse()
	port, ok := config.RoleToPort[*role]
	if !ok {
		panic("invalid role")
	}
	log.Printf("start service %s on port %d", *role, port)

	http.HandleFunc("/formatter", formatHandler)
	http.HandleFunc("/publisher", publishHandler)

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
