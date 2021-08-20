package main

import (
	"fmt"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}
	helloTo := os.Args[1]
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)

	tracer, closer := initJaeger("hello-world")
	defer closer.Close()

	span := tracer.StartSpan("say_hello")
	span.SetTag("hello-to", helloTo)
	println(helloStr)
	span.Finish()
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	return tracer, closer
}
