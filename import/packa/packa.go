package packa

import (
	"fmt"
	"net/http"
)

var register map[string]string

func init() {
	register = make(map[string]string, 0)
	http.HandleFunc("/Index", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	res := fmt.Sprintf("Hello, %s", r.Header.Get("X-Real-Ip"))
	w.Write([]byte(res))
}

func Put(key string, value string) {
	register[key] = value
}

func Get(key string) {
	fmt.Println(register[key])
}
