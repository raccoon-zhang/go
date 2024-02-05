package gptChat

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type LocalClient struct {
	client    *openai.Client
	compleReq openai.ChatCompletionRequest
	mtx       *sync.Mutex
}

var apiKey string
var defaultGptConfig openai.ClientConfig

func init() {
	content, err := os.ReadFile("../apiKey")
	if err != nil {
		fmt.Print("get apiKey err:", err)
	}
	//防止apiKey末尾有换行符
	apiKey = strings.TrimRight(string(content), "\n")
	defaultGptConfig = openai.DefaultConfig(apiKey)
	defaultGptConfig.BaseURL = "https://api.chatanywhere.com.cn/v1"
}

func DefaultClient() LocalClient {
	return LocalClient{
		client: openai.NewClientWithConfig(defaultGptConfig),
		compleReq: openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: make([]openai.ChatCompletionMessage, 0),
		},
		mtx: &sync.Mutex{},
	}
}

func (c *LocalClient) AddMsg(userMsg string) {
	c.mtx.Lock()
	c.compleReq.Messages = append(c.compleReq.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})
	c.mtx.Unlock()
}

func (c LocalClient) QueryGpt(userMsg string) (interface{}, error) {
	c.AddMsg(userMsg)
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		c.compleReq,
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}
	return resp.Choices[0].Message.Content, nil
}
