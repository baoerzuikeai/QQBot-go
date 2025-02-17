package deepseek

import (
	"context"
	"log"
	"os"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

func Chat(msg string) string {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_KEY"))
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleSystem, Content: "用孙吧老哥的语气回复我"}, //用擅长捉弄的高木同学里面的高木同学语气回复我
			{Role: constants.ChatMessageRoleUser, Content: msg},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	return response.Choices[0].Message.Content
}
