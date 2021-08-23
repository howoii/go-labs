package main

import (
	"fmt"
	"net/url"

	conf "github.com/labs/tracing/config"
	"github.com/labs/tracing/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func getHttpResponse(method string, data url.Values, span opentracing.Span) (string, error) {
	serverUrl := util.GetServerUrl(conf.RoleToPort[method], method)
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, serverUrl)
	ext.HTTPMethod.Set(span, method)

	req := util.GetRequest(serverUrl, data)
	if req == nil {
		return "", fmt.Errorf("make request failed")
	}

	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	body, err := util.DoRequest(req)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
