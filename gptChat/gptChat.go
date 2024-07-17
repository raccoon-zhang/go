package gptChat

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type LocalClient struct {
	client    *openai.Client
	compleReq *openai.ChatCompletionRequest
	mtx       *sync.Mutex
}

var apiKey string
var defaultGptConfig openai.ClientConfig

func init() {
	content, err := os.ReadFile("../apiKey")
	if err != nil {
		slog.Error(err.Error())
	}
	//防止apiKey末尾有换行符
	apiKey = strings.TrimRight(string(content), "\n")
	defaultGptConfig = openai.DefaultConfig(apiKey)
	defaultGptConfig.BaseURL = "https://api.chatanywhere.com.cn/v1"
}

func DefaultClient() LocalClient {
	return LocalClient{
		client: openai.NewClientWithConfig(defaultGptConfig),
		compleReq: &openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: make([]openai.ChatCompletionMessage, 0),
		},
		mtx: &sync.Mutex{},
	}
}

func (c *LocalClient) addMsg(info interface{}) {
	c.mtx.Lock()
	if value, ok := info.(string); ok {
		c.compleReq.Messages = append(c.compleReq.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: value,
		})
	}
	if value, ok := info.(openai.ChatCompletionMessage); ok {
		c.compleReq.Messages = append(c.compleReq.Messages, value)
	}
	c.mtx.Unlock()
}

func (c *LocalClient) QueryGpt(userMsg string) (interface{}, error) {
	c.addMsg(userMsg)
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		*c.compleReq,
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	//记录聊天信息，用于连续对话
	//tips:理论上这样挺占内存的，考虑一下本地缓存做一下，但是磁盘存取效率有点低，需要做一下取舍。
	// - 现在就直接放内存了
	c.addMsg(resp.Choices[0].Message)
	return resp.Choices[0].Message.Content, nil
}
