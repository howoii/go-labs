package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labs/tracing/util"
	"github.com/opentracing/opentracing-go"
)

type TracerHolder interface {
	GetTracer() opentracing.Tracer
}

type Service struct {
	Tracer opentracing.Tracer
	Closer io.Closer
}

func NewService(name string) *Service {
	tracer, closer := util.InitTracer(name)
	if tracer == nil {
		return nil
	}
	return &Service{
		Tracer: tracer,
		Closer: closer,
	}
}

func (srv *Service) FormatHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	helloTo := r.FormValue("helloTo")
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)

	_, err := w.Write([]byte(helloStr))
	return err
}

func (srv *Service) PublishHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	helloStr := r.FormValue("helloStr")
	log.Println(helloStr)

	_, err := w.Write([]byte{})
	return err
}

// GetTracer implement: TracerHolder
func (srv *Service) GetTracer() opentracing.Tracer {
	return srv.Tracer
}
