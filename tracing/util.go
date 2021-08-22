package main

import (
	"fmt"
	"io"
	"net/url"

	pkgConfig "github.com/labs/tracing/config"
	httpUtil "github.com/labs/tracing/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

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

func getHttpResponse(method string, data url.Values, span opentracing.Span) (string, error) {
	serverUrl := httpUtil.GetServerUrl(pkgConfig.RoleToPort[method], method)
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, serverUrl)
	ext.HTTPMethod.Set(span, method)

	req := httpUtil.GetRequest(serverUrl, data)
	if req == nil {
		return "", fmt.Errorf("make request failed")
	}

	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	body, err := httpUtil.Do(req)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
