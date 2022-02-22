package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labs/ratelimit/tokenbucket"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("time: %d\n", time.Now().Unix())
	w.Write([]byte("hello world"))
}

func main() {
	limiter := tokenbucket.New(10, 100)
	lm := newLimiterMiddleware(limiter)

	http.Handle("/", lm.Handle(index))
	log.Fatal(http.ListenAndServe(":7070", nil))
}
