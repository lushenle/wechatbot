package gpt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lushenle/wechatbot/config"
	gogpt "github.com/sashabaranov/go-openai"
)

func Completions(msg string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	if msg == "" {
		return "", fmt.Errorf("request messages required")
	}

	request := &gogpt.ChatCompletionRequest{
		Model:            cfg.Model,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             1,
		PresencePenalty:  0.9,
		FrequencyPenalty: 0.9,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    gogpt.ChatMessageRoleSystem,
				Content: "You're an AI assistant, and I need you to simulate a programmer to answer my questions.",
			},
			{
				Role:    gogpt.ChatMessageRoleUser,
				Content: msg,
			},
		},
	}

	gptConfig := gogpt.DefaultConfig(cfg.ApiKey)

	// detect and configuration proxy
	if cfg.Proxy != "" {
		// creates http transport object, sets proxy
		proxyUrl, err := url.Parse(cfg.Proxy)
		if err != nil {
			log.Fatalf("parse proxy err: %v", err)
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		gptConfig.HTTPClient = &http.Client{Transport: transport}
	}

	// create client object
	client := gogpt.NewClientWithConfig(gptConfig)

	if strings.EqualFold(cfg.Model, gogpt.GPT3Dot5Turbo) || strings.EqualFold(cfg.Model, gogpt.GPT3Dot5Turbo0301) {
		resp, err := client.CreateChatCompletion(context.Background(), *request)
		return resp.Choices[0].Message.Content, err
	}

	prompt := ""
	for _, item := range request.Messages {
		prompt += item.Content + "/n"
	}
	prompt = strings.Trim(prompt, "/n")
	log.Printf("request prompt: %v\n", prompt)

	req := gogpt.CompletionRequest{
		Model:            cfg.Model,
		MaxTokens:        cfg.MaxTokens,
		TopP:             1,
		FrequencyPenalty: 0.9,
		PresencePenalty:  0.9,
		Temperature:      0.5,
		Prompt:           prompt,
	}
	resp, err := client.CreateCompletion(context.Background(), req)

	return resp.Choices[0].Text, err
}
