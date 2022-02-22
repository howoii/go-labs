package main

import (
	"fmt"
	"net/http"
)

type Limiter interface {
	Take() bool
}

type LimiterMiddleware struct {
	limiter  Limiter
	failFunc func(w http.ResponseWriter, r *http.Request)
}

type Option func(*LimiterMiddleware)

func WithFailFunc(f func(w http.ResponseWriter, r *http.Request)) Option {
	return func(lm *LimiterMiddleware) {
		lm.failFunc = f
	}
}

func newLimiterMiddleware(limiter Limiter, opts ...Option) *LimiterMiddleware {
	lm := &LimiterMiddleware{
		limiter: limiter,
	}
	for _, opt := range opts {
		opt(lm)
	}
	return lm
}

func (lm *LimiterMiddleware) Handle(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !lm.limiter.Take() {
			if lm.failFunc != nil {
				lm.failFunc(w, r)
				return
			}
			fmt.Printf("From: %s, request refused by rate limiter\n", r.RemoteAddr)
			http.Error(w, "request refused", http.StatusTooManyRequests)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}
