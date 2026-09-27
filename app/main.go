package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
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

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseURL))
	model := os.Getenv("LOCAL_MODEL")
	if model == "" {
		if os.Getenv("LOCAL") == "true" {
			model = "nvidia/nemotron-3-ultra-550b-a55b:free"
		} else {
			model = "anthropic/claude-haiku-4.5"
		}

	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model: model,
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Tools: []openai.ChatCompletionToolUnionParam{
				openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
					Name: "Read",
					Description: openai.String("Read and return the contents of the file"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]any{
							"file_path": map[string]any{
								"type": "string",
								"description": "Path to the file to read",
						},
					},
						"required": []string{"file_path"},
					},

				}),
			},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		fmt.Fprintln(os.Stderr, "error: no choices in response")
		os.Exit(1)
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	fmt.Print(resp.Choices[0].Message.Content)
}
