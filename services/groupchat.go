package services

import (
	"cloopy/utils"
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"log"
	"net/http"
)

type GroupChatResource struct {
	MessageType   string          `json:"msgtype"`
	Text          TextContent     `json:"text,omitempty"`
	Markdown      MarkdownContent `json:"markdown,omitempty"`
	MentionedList []string        `json:"mentioned_list"`
	WebhookUrl    string          `json:"webhook_url,omitempty"`
	News          NewsContent     `json:"news,omitempty"`
}

//var cloopy = make([]Cloopy, 1)

type GrafanaResource struct {
	Title       string                   `json:"title"`
	RuleId      int                      `json:"ruleId"`
	RuleName    string                   `json:"ruleName"`
	RuleUrl     string                   `json:"ruleUrl"`
	State       string                   `json:"state"`
	ImageUrl    string                   `json:"imageUrl"`
	Message     string                   `json:"message"`
	EvalMatches []map[string]interface{} `json:"evalMatches"`
}

func (u GroupChatResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/cloopy").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	tags := []string{"cloopy"}

	ws.Route(ws.POST("/send").To(u.sendMessage).
		Doc("send message via cloopy.").
		Reads(GroupChatResource{}).Metadata(restfulspec.KeyOpenAPITags, tags))

	ws.Route(ws.POST("/grafana").To(u.grafanaMessage).
		Doc("send message from grafana via cloopy.").
		Param(ws.QueryParameter("webhook", "custom webhook url")).
		Reads(GrafanaResource{}).Metadata(restfulspec.KeyOpenAPITags, tags))
	return ws
}

func (u GroupChatResource) grafanaMessage(req *restful.Request, res *restful.Response) {
	grafana := GrafanaResource{}
	err := req.ReadEntity(&grafana)
	if err != nil {
		log.Println(err)
		_ = res.WriteError(http.StatusInternalServerError, err)
	}
	fmt.Println(grafana)
	var content string
	message := grafana.Message
	if message == "" {
		message = "State Change"
	}
	if grafana.State == "ok" {
		content = fmt.Sprintf("### <font color=\"info\">%s</font> \n"+
			">Rule: <font color=\"comment\">%s</font> \n"+
			">Message: <font color=\"comment\">%s</font>",
			grafana.Title, grafana.RuleName, message)
	} else if len(grafana.EvalMatches) == 0 {
		content = fmt.Sprintf("### <font color=\"warning\">%s</font> \n"+
			">Rule: <font color=\"comment\">%s</font> \n"+
			">Message: <font color=\"comment\">%s</font>",
			grafana.Title, grafana.RuleName, message)
	} else {
		content = fmt.Sprintf("### <font color=\"warning\">%s</font> \n"+
			">Rule: <font color=\"comment\">%s</font> \n"+
			">Metrics: <font color=\"comment\">%s</font> \n"+
			">Value: <font color=\"comment\">%s </font> \n"+
			">Message: <font color=\"comment\">%s </font>",
			grafana.Title, grafana.RuleName, grafana.EvalMatches[0]["metric"], grafana.EvalMatches[0]["value"], message)
	}
	msg := GroupChatResource{
		MessageType: "markdown",
		Markdown: MarkdownContent{
			Content: content,
		},
	}
	webhook := req.QueryParameter("webhook")
	body, err := utils.POST(webhook, msg)
	if err != nil {
		_ = res.WriteError(http.StatusInternalServerError, err)
	}
	fmt.Println(string(body))
	_ = res.WriteHeaderAndEntity(http.StatusCreated, grafana)
}

func (u GroupChatResource) sendMessage(req *restful.Request, res *restful.Response) {
	cloopy := GroupChatResource{}
	err := req.ReadEntity(&cloopy)
	if err != nil {
		log.Println(err)
		_ = res.WriteError(http.StatusInternalServerError, err)
	}
	fmt.Println(cloopy)
	var msg GroupChatResource
	msgType := cloopy.MessageType
	if msgType == "text" {
		msg = GroupChatResource{
			MessageType: msgType,
			Text: TextContent{
				Content:       cloopy.Text.Content,
				MentionedList: cloopy.Text.MentionedList,
			},
		}
	} else if msgType == "markdown" {
		msg = GroupChatResource{
			MessageType: msgType,
			Markdown: MarkdownContent{
				Content: cloopy.Markdown.Content,
			},
		}
	} else if msgType == "news" {
		msg = GroupChatResource{
			MessageType: msgType,
			News: NewsContent{
				Articles: cloopy.News.Articles,
			},
		}
	}

	webhookUrl := cloopy.WebhookUrl
	body, err := utils.POST(webhookUrl, msg)
	if err != nil {
		_ = res.WriteError(http.StatusInternalServerError, err)
	}
	fmt.Println(string(body))
	_ = res.WriteHeaderAndEntity(http.StatusCreated, cloopy)
}

type TextContent struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
}

type MarkdownContent struct {
	Content string `json:"content"`
}

type NewsContent struct {
	Articles []map[string]interface{} `json:"articles"`
}
