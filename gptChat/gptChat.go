package gptChat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

//url及参数信息详见：https://github.com/bluoruo/go-gtp3.5

type msgBox struct {
	Role    string `json:"role"` //"user"代表话题内容，"system"代表话题类型
	Content string `json:"content"`
}

// POST方法参数列表
type RequestInfo struct {
	Messages         []msgBox `json:"messages"`
	Model            string   `json:"model"`
	MaxTokens        uint     `json:"max_tokens"`
	Temperature      float64  `json:"temperature"`
	TopP             uint     `json:"top_p"`
	Stream           bool     `json:"stream"`
	N                uint     `json:"n"`
	PresencePenalty  float64  `json:"presence_penalty"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	User             string   `json:"user"`
}

type UsageInfo struct {
	PromptTokens     uint `json:"prompt_tokens"`
	CompletionTokens uint `json:"completion_tokens"`
	TotalTokens      uint `json:"total_tokens"`
}

// gpt返回信息
type responceInfo struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created uint      `json:"created"`
	Usage   UsageInfo `json:"usage"`
	Choices struct {
		Msg          msgBox `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        uint   `json:"index"`
	}
}

var gptUrl = "https://api.openai.com/v1/chat/completions"

func (defaultReqInfo *RequestInfo) setDefatuReqInfo() {
	defaultReqInfo.Model = "gpt-3.5-turbo"
	defaultReqInfo.PresencePenalty = 1.0
	defaultReqInfo.FrequencyPenalty = 1.0
	defaultReqInfo.User = "default"
	var msg = msgBox{
		Role:    "system",
		Content: "developer",
	}
	defaultReqInfo.Messages = make([]msgBox, 0)
	defaultReqInfo.Messages = append(defaultReqInfo.Messages, msg)
}

func QueryGpt(userMsg string) (interface{}, error) {
	var reqInfo RequestInfo
	reqInfo.setDefatuReqInfo()
	reqInfo.Messages = append(reqInfo.Messages, msgBox{
		Role:    "user",
		Content: userMsg,
	})
	jsonData, err := json.Marshal(reqInfo)
	if err != nil {
		return nil, err
	}
	responce, err := http.Post(gptUrl, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	} else {
		defer responce.Body.Close()
		responceBody, err := io.ReadAll(responce.Body)
		if err != nil || responce.StatusCode != http.StatusOK {
			fmt.Println("err:", err)
			fmt.Println("responceBody:", responceBody)
			return nil, errors.New("read responceBody wrong")
		} else {
			var responceJsonData responceInfo
			err := json.Unmarshal(responceBody, &responceJsonData)
			if err != nil {
				return nil, err
			}
			return responceJsonData.Choices.Msg, nil
		}
	}
}
