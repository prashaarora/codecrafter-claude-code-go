package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	_ = godotenv.Load()
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		fmt.Fprintln(os.Stderr, "error: prompt must not be empty")
		flag.Usage()
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: env variable OPENROUTER_API_KEY not found")
		os.Exit(1)
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	model := os.Getenv("LOCAL_MODEL")
	if model == "" {
		model = "anthropic/claude-haiku-4.5"
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}

	tools := buildTools()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := client.Chat.Completions.New(ctx,
			openai.ChatCompletionNewParams{
				Model:    model,
				Messages: messages,
				Tools:    tools,
			},
		)
		cancel()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(resp.Choices) == 0 {
			fmt.Fprintln(os.Stderr, "error: no choices in response")
			os.Exit(1)
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg.ToParam())
		if len(msg.ToolCalls) == 0 {
			fmt.Print(msg.Content)
			return
		}

		for _, toolCall := range msg.ToolCalls {
			result, err := executeToolCall(toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to execute tool call: %v\n", err)
				os.Exit(1)
			}
			messages = append(messages, openai.ToolMessage(result, toolCall.ID))
		}
	}
}
