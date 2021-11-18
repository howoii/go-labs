package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var (
	appID  = "********"
	appKey = "********"

	fromLang = "en"
	toLang   = "zh"

	endpoint   = "http://api.fanyi.baidu.com"
	path       = "/api/trans/vip/translate"
	requestUrl = endpoint + path

	query = flag.String("query", "", "source text to translate")
)

func makeMD5(txt string) string {
	h := md5.New()
	h.Write([]byte(txt))
	return hex.EncodeToString(h.Sum(nil))
}

// random in [begin, end]
func randBetween(begin int, end int) int {
	if end < begin {
		return rand.Intn(begin-end+1) + end
	}
	return rand.Intn(end-begin+1) + begin
}

type transResult struct {
	Src string `json:"src,omitempty"`
	Dst string `json:"dst,omitempty"`
}

type transResponse struct {
	From        string        `json:"from,omitempty"`
	To          string        `json:"to,omitempty"`
	TransResult []transResult `json:"trans_result"`
}

func main() {
	flag.Parse()
	if *query == "" {
		log.Println("you need to input translation text")
		return
	}
	salt := strconv.FormatInt(int64(randBetween(32768, 65536)), 10)
	sign := makeMD5(appID + *query + salt + appKey)

	values := url.Values{}
	values.Set("appid", appID)
	values.Set("q", *query)
	values.Set("from", fromLang)
	values.Set("to", toLang)
	values.Set("salt", salt)
	values.Set("sign", sign)
	resp, err := http.PostForm(requestUrl, values)
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var res transResponse
	if err := json.Unmarshal(body, &res); err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Printf("en:\n%s\nzh:\n%s\n", res.TransResult[0].Src, res.TransResult[0].Dst)
}
