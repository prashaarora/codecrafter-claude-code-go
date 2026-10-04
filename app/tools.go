package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go/v3"
)

type ReadArgs struct {
	FilePath string `json:"file_path"`
}
type WriteArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type BashArgs struct {
	Command string `json:"command"`
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
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "Write",
			Description: openai.String("Write content to a file"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "Path to the file to write to",
					},
					"content": map[string]any{
						"type":        "string",
						"description": "Content to write to the file",
					},
				},
				"required": []string{"file_path", "content"},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "Bash",
			Description: openai.String("Execute a shell command"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
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
	case "Write":
		return executeWriteTool(arguments)
	case "Bash":
		return executeBashTool(arguments)
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

func executeWriteTool(arguments string) (string, error) {
	var args WriteArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}
	if err := os.WriteFile(args.FilePath, []byte(args.Content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return "success", nil
}

func executeBashTool(arguments string) (string, error) {
	var args BashArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}
	cmd := exec.Command("sh", "-c", args.Command)
	output, err := cmd.CombinedOutput()
	result := string(output)
	if err != nil {
		if result == "" {
			result = err.Error()
		} else {
			result = result + "\n" + err.Error()
		}
		return result, err
	}
	if result == "" {
		result = "Command executed successfully"
	}
	return result, nil
}
