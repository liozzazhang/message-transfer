package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/micro/go-micro/config/encoder/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	webhook := flag.String("webhook_url", "cloopy", "webhook url")
	msgType := flag.String("msg_type", "text", "Default is text, `text|markdown` available.")
	message := flag.String("message", "Hello, This is Cloopy", "message contents")
	mentionedListString := flag.String("mention", "", "mention people, Comma delimited list of project name")
	flag.Parse()
	mentionedList := strings.Split(*mentionedListString, ",")

	var msg webhookResource
	if *msgType == "text" {
		msg = webhookResource{
			MessageType: *msgType,
			Text: textContent{
				Content:       *message,
				MentionedList: mentionedList,
			},
		}
	} else if *msgType == "markdown" {
		msg = webhookResource{
			MessageType: *msgType,
			Markdown: markdownContent{
				Content:       *message,
				MentionedList: mentionedList,
			},
		}
	}

	if *webhook == "cloopy" {
		*webhook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3b179ce0-9da0-49b3-9847-88f9d2ae142f"
	}

	Url, err := url.Parse(*webhook)
	if err != nil {
		panic(err)
	}
	urlPath := Url.String()
	//buf := new(bytes.Buffer)
	buf, _ := json.NewEncoder().Encode(msg)
	fmt.Println(string(buf))
	request, _ := http.NewRequest("POST", urlPath, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

type webhookResource struct {
	MessageType string          `json:"msgtype"`
	Text        textContent     `json:"text,omitempty"`
	Markdown    markdownContent `json:"markdown,omitempty"`
	//MentionedMobileList []string `json:"mentioned_mobile_list"`
}

type textContent struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}

type markdownContent struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}
