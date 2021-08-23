package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	conf "github.com/labs/tracing/config"
)

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
	resp, err := http.DefaultClient.Do(req)
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
