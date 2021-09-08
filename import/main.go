// this lab is used to test whether go import each package once
package main

import (
	"net/http"

	"github.com/labs/import/packa"
	"github.com/labs/import/packb"
)

func main() {
	packb.Register()
	packa.Put("main", "this is package main")
	packa.Get("packb")

	http.HandleFunc("/Hello", helloHandler)
	if err := http.ListenAndServe(":8099", nil); err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
