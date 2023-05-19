package services

import (
	"context"
	"fmt"
	"gptube/config"

	"github.com/sashabaranov/go-openai"
)

var ChatGPTClient = openai.NewClient(config.Config("OPENAI_API_KEY"))

func Chat(message string) (*openai.ChatCompletionResponse, error) {
	resp, err := ChatGPTClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	return &resp, nil
}
