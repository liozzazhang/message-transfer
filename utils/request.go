package utils

import (
	"bytes"
	"fmt"
	"github.com/micro/go-micro/config/encoder/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func POST(httpUrl string, msg interface{}) (string, error) {
	if httpUrl == "" {
		httpUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3b179ce0-9da0-49b3-9847-88f9d2ae142f"
	}

	Url, err := url.Parse(httpUrl)
	if err != nil {
		return "url parse error", err
	}
	urlPath := Url.String()
	//buf := new(bytes.Buffer)
	buf, _ := json.NewEncoder().Encode(msg)
	fmt.Println(string(buf))
	request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}
