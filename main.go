package main

import (
	"encoding/json"
	"github.com/eawsy/aws-lambda-go/service/lambda/runtime"
	"strings"
	"encoding/hex"
	"fmt"
)

type Response struct {
	StatusCode string `json:"statusCode"`
	Body Message `json:"body"`
	Headers map[string]string `json:"headers"`
}

//Type Message
type Message struct{
	ResponseType string `json:"response_type"`
	Text string `json:"text"`
}
//Type Request
type Request struct {
	Resource string `json:"resource"`
	Path string `json:"path"`
	HttpMethod string `json:"httpMethod"`
	Headers map[string]string `json:"headers"`
	QueryStringParameters string `json:"queryStringParameters"`
	PathParameters string `json:"pathParameters"`
	StageVariables string `json:"stageVariables"`
	RequestContext map[string]string `json:"requestContext"`
	Body string `json:"body"`
	IsBase64Encoded string `json:"isBase64Encoded"`
}

type Info struct{
	TeamId string `json:"team_id"`
	TeamDomain string`json:"team_domain"`
	ChannelId string`json:"channel_id"`
	UserId string `json:"user_id"`
	UserName string`json:"user_name"`
	Command string`json:"command"`
	ResponseUrl string `json:"response_url"`
	Token string `json:"token"`
	ChannelName string`json:"channel_name"`
	Text string`json:"text"`
}

func toInfo (input string) Info {
	list := strings.Split(input, "&")
	data := make(map[string]string)
	for _, x := range list {
		temp := strings.Split(x, "=")
		data[temp[0]] = temp[1]
	}
	for key, str := range data {
		changes := make(map[string]string)
		for index, value := range str {
			if value == '%' {
				decoded, _ := hex.DecodeString(str[index+1 : index+3])
				changes[str[index:index+3]] = string(decoded)
			}
		}
		for original, change := range changes {
			data[key] = strings.Replace(data[key], original, change, -1)
		}
	}
	result := Info{
		TeamId: data["team_id"],
		TeamDomain: data["team_domain"],
		ChannelId: data["channel_id"],
		UserId: data["user_id"],
		UserName: data["user_name"],
		Command: data["command"],
		ResponseUrl: data["response_url"],
		Token: data["token"],
		ChannelName: data["channel_name"],
		Text: data["text"]}
	return result
}

func handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	var request Request
	json.Unmarshal(evt, &request)
	info := toInfo(request.Body)
	headers := map[string]string{"Content-Type": "application/json"}
	message := Message{ResponseType:"ephemeral",Text:info.Text}
	result := Response{StatusCode: "200", Body:message, Headers:headers}

	fmt.Println("sending response")
	return result, nil

}

func init() {
	runtime.HandleFunc(handle)
}

func main() {}
