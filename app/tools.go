package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

type ReadArgs struct {
	FilePath string `json:"file_path"`
}

func buildTools() []openai.ChatCompletionToolUnionParam {
	return []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "Read",
			Description: openai.String("Read and return the contents of the file"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "Path to the file to read",
					},
				},
				"required": []string{"file_path"},
			},
		}),
	}
}

func executeToolCall(toolName string, arguments string) (string, error) {
	return executeTool(toolName, arguments)
}

func executeTool(toolName string, arguments string) (string, error) {
	switch toolName {
	case "Read":
		return executeReadTool(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

func executeReadTool(arguments string) (string, error) {
	var args ReadArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}
	content, err := os.ReadFile(args.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}
