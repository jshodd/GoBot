package main

import (
	"encoding/json"
	"fmt"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"net/url"
	"strings"
)

type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

func NewResponse(statusCode int, responseType, text string) Response {
	headers := map[string]string{"Content-Type": "application/json"}
	message := Message{ResponseType: responseType, Text: text}
	return Response{StatusCode: statusCode, Body: message.String(), Headers: headers}
}

//Type Message
type Message struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func (this *Message) String() string {
	return fmt.Sprintf("{'response_type': '%s', 'text': \"%s\"}", this.ResponseType, this.Text)
}

//Type Request
type Request struct {
	Resource              string            `json:"resource"`
	Path                  string            `json:"path"`
	HttpMethod            string            `json:"httpMethod"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters string            `json:"queryStringParameters"`
	PathParameters        string            `json:"pathParameters"`
	StageVariables        string            `json:"stageVariables"`
	RequestContext        map[string]string `json:"requestContext"`
	Body                  string            `json:"body"`
	IsBase64Encoded       string            `json:"isBase64Encoded"`
}

type SlackInfo struct {
	TeamId      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelId   string `json:"channel_id"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	ResponseUrl string `json:"response_url"`
	Token       string `json:"token"`
	ChannelName string `json:"channel_name"`
	Text        string `json:"text"`
}

func toInfo(input string) SlackInfo {
	data, _ := url.ParseQuery(input)
	result := SlackInfo{
		TeamId:      data["team_id"][0],
		TeamDomain:  data["team_domain"][0],
		ChannelId:   data["channel_id"][0],
		UserId:      data["user_id"][0],
		UserName:    data["user_name"][0],
		Command:     data["command"][0],
		ResponseUrl: data["response_url"][0],
		Token:       data["token"][0],
		ChannelName: data["channel_name"][0],
		Text:        data["text"][0]}
	return result
}

func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	var request Request
	json.Unmarshal(evt, &request)
	info := toInfo(request.Body)
	if info.Command == "/infobot" {
		switch info.Text {
		case "oscprules":
			text := "1. No spoilers\n2. No irrelevant discussions/offtopic chatter (okay in threads)\n3. Research before asking (http://www.catb.org/esr/faqs/smart-questions.html#before)\n4. NO DIRECT QUESTIONS ON MACHINES\n5. No shitposting/memes\n6. Try harder 7.\nPeople blatantly ignoring the rules will either get warned or kicked out from the channel.\n8. Don't add members to the channel, instead ask an admin."
			return NewResponse(200, "ephemeral", text), nil
		case "learn":
			text := "NetSec Learning Resources: https://docs.google.com/spreadsheets/d/12bT8APhWsL-P8mBtWCYu4MLftwG1cPmIL25AEBtXDno/edit?usp=sharing"
			return NewResponse(200, "ephemeral", text), nil
		default:
			text := "I don't know that one yet, this is just a proof of concept, but I do know:\n1) `/infobot learn` \n2) `/infobot oscprules`"
			return NewResponse(200, "ephemeral", text), nil
		}
	} else if info.Command == "/lmgtfy" {
		text := "http://lmgtfy.com/?q=" + strings.Replace(info.Text, " ", "+", -1)
		return NewResponse(200, "in_channel", text), nil
	}
	return NewResponse(500, "ephemeral", "something went wrong..."), nil

}
