package gptChat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

//url及参数信息详见：https://github.com/bluoruo/go-gtp3.5

// type msgBox struct {
// 	Role    string `json:"role"` //"user"代表话题内容，"system"代表话题类型
// 	Content string `json:"content"`
// }

// // POST方法参数列表
// type RequestInfo struct {
// 	Messages         []msgBox `json:"messages"`
// 	Model            string   `json:"model"`
// 	MaxTokens        uint     `json:"max_tokens"`
// 	Temperature      float64  `json:"temperature"`
// 	TopP             uint     `json:"top_p"`
// 	Stream           bool     `json:"stream"`
// 	N                uint     `json:"n"`
// 	PresencePenalty  float64  `json:"presence_penalty"`
// 	FrequencyPenalty float64  `json:"frequency_penalty"`
// 	User             string   `json:"user"`
// }

// type UsageInfo struct {
// 	PromptTokens     uint `json:"prompt_tokens"`
// 	CompletionTokens uint `json:"completion_tokens"`
// 	TotalTokens      uint `json:"total_tokens"`
// }

// // gpt返回信息
// type responceInfo struct {
// 	Id      string    `json:"id"`
// 	Object  string    `json:"object"`
// 	Created uint      `json:"created"`
// 	Usage   UsageInfo `json:"usage"`
// 	Choices struct {
// 		Msg          msgBox `json:"message"`
// 		FinishReason string `json:"finish_reason"`
// 		Index        uint   `json:"index"`
// 	}
// }

var gptUrl string
var apiKey string
var defaultCompleReq openai.ChatCompletionRequest

func init() {
	defaultCompleReq = openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: make([]openai.ChatCompletionMessage, 0),
	}
	gptUrl = "https://api.chatanywhere.com.cn/v1"
	content, err := os.ReadFile("../apiKey")
	if err != nil {
		fmt.Print("get apiKey err:", err)
	}
	//防止apiKey末尾有换行符
	apiKey = strings.TrimRight(string(content), "\n")
}

// func (defaultReqInfo *RequestInfo) setDefatuReqInfo() {
// 	defaultReqInfo.Model = "gpt-3.5-turbo"
// 	defaultReqInfo.PresencePenalty = 1.0
// 	defaultReqInfo.FrequencyPenalty = 1.0
// 	defaultReqInfo.User = "default"
// 	var msg = msgBox{
// 		Role:    "system",
// 		Content: "developer",
// 	}
// 	defaultReqInfo.Messages = make([]msgBox, 0)
// 	defaultReqInfo.Messages = append(defaultReqInfo.Messages, msg)
// }

func QueryGpt(userMsg string) (interface{}, error) {
	client := openai.NewClient(apiKey)
	defaultCompleReq.Messages = append(defaultCompleReq.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})
	resp, err := client.CreateChatCompletion(
		context.Background(),
		defaultCompleReq,
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	return resp.Choices[0].Message.Content, nil
}
