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
			Model: openai.GPT3Dot5Turbo0301,
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

func (ct *OpenAI) GetImage(desc string) (string, error) {

	resp, err := ct.OpenAI.CreateImage(context.Background(), openai.ImageRequest{
		Prompt: desc,
		N:      1,
		Size:   "512x512",
	})
	if err != nil {
		return "", err
	}
	return resp.Data[0].URL, nil
}
