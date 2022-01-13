package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	conf "github.com/labs/tracing/config"
)

var uTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var uClient = &http.Client{
	Transport: uTransport,
}

func GetServerUrl(port int32, method string) string {
	return fmt.Sprintf(conf.ServerUrl, port, method)
}

func GetRequest(url string, body url.Values) *http.Request {
	req, err := http.NewRequest("POST", url, strings.NewReader(body.Encode()))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func DoRequest(req *http.Request) ([]byte, error) {
	resp, err := uClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status not ok")
	}

	return body, nil
}
