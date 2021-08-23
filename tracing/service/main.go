package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	conf "github.com/labs/tracing/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

	srv := NewService(*method)
	defer srv.Closer.Close()

	http.HandleFunc("/format", handlerWrap(srv, srv.FormatHandler))
	http.HandleFunc("/publish", handlerWrap(srv, srv.PublishHandler))

	url := ":" + strconv.FormatInt(int64(port), 10)
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Println("server exit")
	}
}

type httpMethodFunc func(w http.ResponseWriter, r *http.Request) error

func handlerWrap(th TracerHolder, handler httpMethodFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracer := th.GetTracer()
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		method := r.URL.Path
		var span opentracing.Span
		if err != nil {
			span = tracer.StartSpan(method)
		} else {
			span = tracer.StartSpan(method, ext.RPCServerOption(spanCtx))
		}
		defer span.Finish()

		err = handler(w, r)
		if err != nil {
			ext.LogError(span, err)
		}
	}
}
