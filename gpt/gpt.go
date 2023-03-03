package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lushenle/wechatbot/config"
	"github.com/lushenle/wechatbot/pkg/logger"
)

//const BASEURL = "https://api.openai.com/v1/chat/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

// GPT Response
//	{
//	 "id": "chatcmpl-123",
//	 "object": "chat.completion",
//	 "created": 1677652288,
//	 "choices": [{
//	   "index": 0,
//	   "message": {
//	     "role": "assistant",
//	     "content": "\n\nHello there, how may I assist you today?",
//	   },
//	   "finish_reason": "stop"
//	 }],
//	 "usage": {
//	   "prompt_tokens": 9,
//	   "completion_tokens": 12,
//	   "total_tokens": 21
//	 }
//	}

type ChoiceItem struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string    `json:"model"`
	Prompt           []Message `json:"messages"`
	MaxTokens        uint      `json:"max_tokens"`
	Temperature      float64   `json:"temperature"`
	TopP             int       `json:"top_p"`
	FrequencyPenalty int       `json:"frequency_penalty"`
	PresencePenalty  int       `json:"presence_penalty"`
}

// Completions GPT request
// https://platform.openai.com/docs/api-reference/chat/create?lang=curl
//
//	curl https://api.openai.com/v1/chat/completions \
//	 -H 'Content-Type: application/json' \
//	 -H 'Authorization: Bearer YOUR_API_KEY' \
//	 -d '{
//	 "model": "gpt-3.5-turbo",
//	 "messages": [{"role": "user", "content": "Hello!"}]
//	}'
//
// Parameters
//
//	{
//	 "model": "gpt-3.5-turbo",
//	 "messages": [{"role": "user", "content": "Hello!"}]
//	}
func Completions(msg string) (string, error) {
	cfg := config.LoadConfig()
	requestBody := ChatGPTRequestBody{
		Model: cfg.Model,
		Prompt: []Message{
			{Role: "user", Content: msg},
		},
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	req, err := http.NewRequest("POST", cfg.API, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Setting the proxy and timeout duration
	//proxy, _ := url.Parse("http://host:port")
	//proxy, _ := url.Parse(cfg.Proxy)
	//client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("请求GTP出错了，gpt api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Message.Content
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	return reply, nil
}
