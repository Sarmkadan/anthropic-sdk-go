package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AdvancedUsage demonstrates advanced configuration, error handling, and best practices
// This example shows how to configure the client with various options and handle errors properly
func main() {
	// Example 1: Using environment variable for API key (recommended for production)
	// Set ANTHROPIC_API_KEY environment variable instead of hardcoding
	// export ANTHROPIC_API_KEY="your-api-key-here"

	// Initialize client with custom configuration
	client := anthropic.NewClient(
		// Use environment variable for API key
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),

		// Custom base URL (useful for testing or self-hosted instances)
		option.WithBaseURL("https://api.anthropic.com"),

		// Set custom timeout for requests
		option.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),

		// Add custom headers (useful for tracking or debugging)
		option.WithHeader("X-Custom-Header", "my-app-v1.0"),
	)

	// Example 2: Using different authentication methods
	// You can also use:
	// - option.WithAuthToken() for bearer tokens
	// - option.WithConfig() for profile-based configuration
	// - option.WithoutEnvironmentDefaults() to disable automatic credential loading

	// Example 3: Making requests with error handling
	ctx := context.Background()

	// Create a conversation with multiple messages
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("Explain quantum computing in simple terms.")),
	}

	// Add previous conversation history if needed
	// messages = append(messages, previousMessages...)

	// Make the API call with timeout
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel() // Always cancel to release resources

	// Create the message request with structured parameters
	requestParams := anthropic.MessageNewParams{
		MaxTokens: 2048,
		Messages:  messages,
		Model:     anthropic.ModelClaudeOpus4_8,

		// Optional parameters for better control
		Temperature: anthropic.FloatPtr(0.7), // 0.0 to 1.0, higher is more creative
		TopP:       anthropic.FloatPtr(0.9),   // Nucleus sampling parameter
		TopK:       anthropic.IntPtr(50),     // Top-k sampling parameter

		// Stop sequences to end generation early
		StopSequences: []string{"\n\n", "Human:", "Assistant:"},
	}

	// Execute the request
	message, err := client.Messages.New(ctx, requestParams)
	if err != nil {
		// Handle different types of errors
		switch err := err.(type) {
		case *anthropic.Error:
			// API returned an error response
			log.Printf("API Error (Status %d): %s\n", err.StatusCode, err.Message)
			log.Printf("Error Type: %s\n", err.Type)
			if err.Error() != nil {
				log.Printf("Details: %v\n", err.Error())
			}
		case *anthropic.APIError:
			// Authentication or validation error
			log.Printf("Authentication/Validation Error: %s\n", err.Message)
		default:
			// Network or other errors
			log.Printf("Request failed: %v\n", err)
		}
		os.Exit(1)
	}

	// Process the response
	fmt.Println("=== Response ===")
	for i, content := range message.Content {
		fmt.Printf("Content %d:\n", i+1)
		switch {
		case content.Text != nil:
			fmt.Printf("Text: %s\n", *content.Text)
		case content.ToolUse != nil:
			fmt.Printf("Tool Use: %s (ID: %s)\n",
				content.ToolUse.Name,
				content.ToolUse.Id,
			)
		case content.Redacted != nil:
			fmt.Printf("Redacted content (block type: %s)\n", content.Type)
		default:
			fmt.Printf("Unknown content type: %s\n", content.Type)
		}
	}

	// Example 4: Working with tool use (function calling)
	if len(message.Content) > 0 && message.Content[0].ToolUse != nil {
		fmt.Println("\n=== Tool Use Detected ===")
		toolUse := message.Content[0].ToolUse

		// Handle the tool call
		result, err := handleToolCall(toolUse)
		if err != nil {
			log.Printf("Error handling tool call: %v\n", err)
			os.Exit(1)
		}

		// Send the tool result back
		toolResultParams := anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{
				// Include the original user message
				anthropic.NewUserMessage(anthropic.NewTextBlock("Explain quantum computing")),
				// Include the assistant's tool use message
				anthropic.NewAssistantMessage(anthropic.NewTextBlock(""),
				anthropic.NewToolUseBlock(toolUse.Id, toolUse.Name, result),
			},
			Model: anthropic.ModelClaudeSonnet4_5_20250929,
		}

		_, err = client.Messages.New(ctx, toolResultParams)
		if err != nil {
			log.Printf("Error sending tool result: %v\n", err)
			os.Exit(1)
		}
	}

	// Example 5: Usage tracking
	fmt.Printf("\n=== Usage Statistics ===")
	fmt.Printf("Total tokens: %d\n", message.Usage.TotalTokens)
	fmt.Printf("Input tokens: %d\n", message.Usage.InputTokens)
	fmt.Printf("Output tokens: %d\n", message.Usage.OutputTokens)
	fmt.Printf("Cache read tokens: %d\n", message.Usage.CacheReadInputTokens)
	fmt.Printf("Cache write tokens: %d\n", message.Usage.CacheWriteInputTokens)
	fmt.Printf("Duration: %v\n", message.Usage.UsageType.Duration)
}

// handleToolCall demonstrates how to implement tool/function calling
func handleToolCall(toolUse *anthropic.ToolUseBlock) (string, error) {
	switch toolUse.Name {
	case "get_current_time":
		return time.Now().Format(time.RFC3339), nil
	case "get_weather":
		return "Sunny, 72°F", nil
	default:
		return "", fmt.Errorf("unknown tool: %s", toolUse.Name)
	}
}
