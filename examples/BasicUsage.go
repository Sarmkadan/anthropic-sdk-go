package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BasicUsage demonstrates the minimal setup required to use the Anthropic SDK
// This example shows how to create a client and make your first API call
func main() {
	// Initialize the client with your API key
	// The API key can be provided directly or via ANTHROPIC_API_KEY environment variable
	client := anthropic.NewClient(
		option.WithAPIKey("your-api-key-here"), // Replace with your actual API key
	)

	// Create a simple message request
	// This sends a text message to the Claude model and gets a response
	message, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024, // Maximum number of tokens to generate
		Messages: []anthropic.MessageParam{
			// Create a user message with text content
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello! Who are you?")),
		},
		Model: anthropic.ModelClaudeSonnet4_5_20250929, // Use an appropriate model
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making API call: %v\n", err)
		os.Exit(1)
	}

	// Print the response content
	fmt.Println("Response:")
	for _, content := range message.Content {
		if content.Text != nil {
			fmt.Println(*content.Text)
		}
	}

	// Print usage information
	fmt.Printf("\nUsage: %d tokens, Input: %d tokens, Output: %d tokens\n",
		message.Usage.TotalTokens,
		message.Usage.InputTokens,
		message.Usage.OutputTokens,
	)
}
