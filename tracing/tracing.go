package main

import (
	"context"
	"net/url"
	"os"

	conf "github.com/labs/tracing/config"
	"github.com/labs/tracing/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	serverUrl = "http://localhost:%d/%s"
)

func main() {
	if len(os.Args) != 3 {
		panic("ERROR: Expecting two argument")
	}
	greeting := os.Args[1]
	helloTo := os.Args[2]

	tracer, closer := util.InitTracer("hello-world")
	if tracer == nil {
		panic("ERROR: init tracer failed")
	}
	opentracing.SetGlobalTracer(tracer) //设置全局tracer，StartSpanFromContext里会用到
	defer closer.Close()

	span := tracer.StartSpan("say_hello")
	span.SetTag("hello-to", helloTo)
	span.SetBaggageItem("greeting", greeting)

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)

	span.Finish()
}

func formatString(ctx context.Context, helloTo string) string {
	// 从ctx中取出rootSpan，开始一个新的span，作为rootSpan的子span（详情看源码注释）
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	resp, err := getHttpResponse(conf.MethodFormat, url.Values{
		"helloTo": {helloTo},
	}, span)
	if err != nil {
		ext.LogError(span, err)
		return resp
	}
	span.LogFields(
		log.String("event", "http-format-string"),
		log.String("value", resp),
	)
	return resp
}

func printHello(ctx context.Context, helloStr string) {
	//原始写法，可以用StartSpanFromContext一步到位
	rootSpan := opentracing.SpanFromContext(ctx)
	span := rootSpan.Tracer().StartSpan("printHello", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()

	_, err := getHttpResponse(conf.MethodPublish, url.Values{
		"helloStr": {helloStr},
	}, span)
	if err != nil {
		ext.LogError(span, err)
		return
	}
	span.LogKV(
		"event", "http-print-hello",
	)
}
