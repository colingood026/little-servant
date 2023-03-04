package main

import (
	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
)

type OpenAI struct {
	OpenAI *openai.Client
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		OpenAI: openai.NewClient(apiKey),
	}
}

func (ct *OpenAI) ChatWithChatGPT(question string) (string, error) {
	resp, err := ct.OpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
