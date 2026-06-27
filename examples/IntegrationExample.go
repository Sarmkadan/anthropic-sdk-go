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
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// IntegrationExample demonstrates how to integrate the Anthropic SDK into a Go application
// This example shows best practices for:
// - Dependency injection patterns
// - Configuration management
// - Error handling middleware
// - Context propagation
// - Graceful shutdown
func main() {
	// Initialize the application with configuration
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// AppConfig holds application configuration
// In a real application, this would come from environment variables or config files
type AppConfig struct {
	APIKey      string
	BaseURL      string
	Timeout      time.Duration
	MaxRetries   int
	Model        constant.Model
	MaxTokens    int
	Temperature  float64
}

// AnthropicService wraps the Anthropic client for easier testing and mocking
type AnthropicService struct {
	client *anthropic.Client
}

// NewAnthropicService creates a new Anthropic service
func NewAnthropicService(cfg AppConfig) (*AnthropicService, error) {
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithHTTPClient(&http.Client{
			Timeout: cfg.Timeout,
		}),
	)

	return &AnthropicService{client: client}, nil
}

// GenerateResponse makes an API call to generate a response from Claude
func (s *AnthropicService) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	// Create the message request
	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		Model:        constant.ModelClaudeSonnet4_5_20250929,
		Temperature:  anthropic.FloatPtr(0.7),
		TopP:        anthropic.FloatPtr(0.9),
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract text content from response
	var responseText string
	for _, content := range message.Content {
		if content.Text != nil {
			responseText = *content.Text
			break
		}
	}

	return responseText, nil
}

// ChatSession represents a conversation session
// This demonstrates stateful conversation management
type ChatSession struct {
	service    *AnthropicService
	messages   []anthropic.MessageParam
	maxHistory int
}

// NewChatSession creates a new chat session
func NewChatSession(service *AnthropicService, maxHistory int) *ChatSession {
	return &ChatSession{
		service:    service,
		messages:   make([]anthropic.MessageParam, 0),
		maxHistory: maxHistory,
	}
}

// SendMessage adds a message to the session and gets a response
func (s *ChatSession) SendMessage(ctx context.Context, userMessage string) (string, error) {
	// Add user message to history
	userMsg := anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage))
	s.messages = append(s.messages, userMsg)

	// Trim history if it exceeds max
	if len(s.messages) > s.maxHistory {
		s.messages = s.messages[len(s.messages)-s.maxHistory:]
	}

	// Get response from Claude
	response, err := s.service.GenerateResponse(ctx, userMessage)
	if err != nil {
		return "", err
	}

	// Add assistant response to history
	assistantMsg := anthropic.NewAssistantMessage(anthropic.NewTextBlock(response), nil)
	s.messages = append(s.messages, assistantMsg)

	return response, nil
}

// run demonstrates the full application lifecycle
func run() error {
	// Load configuration from environment variables
	cfg := loadConfig()

	// Initialize services
	service, err := NewAnthropicService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize Anthropic service: %w", err)
	}

	// Create chat session with 10 message history limit
	session := NewChatSession(service, 10)

	// Example 1: Simple one-off request
	fmt.Println("=== Example 1: Simple Request ===")
	response, err := session.SendMessage(context.Background(), "Hello! Who are you?")
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	fmt.Println("Response:", response)

	// Example 2: Conversation with context
	fmt.Println("\n=== Example 2: Conversation ===")
	ctx := context.Background()

	// Set a timeout for the entire conversation
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Ask a follow-up question
	followUp := "What can you tell me about artificial intelligence?"
	response, err = session.SendMessage(ctx, followUp)
	if err != nil {
		return fmt.Errorf("follow-up request failed: %w", err)
	}
	fmt.Println("Response:", response)

	// Example 3: Using different models for different tasks
	fmt.Println("\n=== Example 3: Model Selection ===")
	// For creative tasks, use a more creative model
	cfg.Model = constant.ModelClaudeOpus4_8
	service, _ = NewAnthropicService(cfg)

	creativePrompt := "Write a short poem about autumn leaves"
	response, err = service.GenerateResponse(ctx, creativePrompt)
	if err != nil {
		return fmt.Errorf("creative request failed: %w", err)
	}
	fmt.Println("Creative Response:", response)

	// Example 4: Error handling and retries
	fmt.Println("\n=== Example 4: Error Handling ===")
	// Simulate handling API errors gracefully
	for i := 0; i < 3; i++ {
		response, err = session.SendMessage(ctx, "Tell me a joke")
		if err != nil {
			log.Printf("Attempt %d failed: %v\n", i+1, err)
			if i < 2 {
				// Wait before retrying
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}
			return fmt.Errorf("all retry attempts failed: %w", err)
		}
		break
	}
	fmt.Println("Joke:", response)

	// Example 5: Usage metrics collection
	fmt.Println("\n=== Example 5: Usage Tracking ===")
	logUsageMetrics(session)

	return nil
}

// loadConfig loads configuration from environment variables
func loadConfig() AppConfig {
	// Default configuration
	cfg := AppConfig{
		APIKey:      os.Getenv("ANTHROPIC_API_KEY"),
		BaseURL:      "https://api.anthropic.com",
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		Model:        constant.ModelClaudeSonnet4_5_20250929,
		MaxTokens:    2048,
		Temperature:  0.7,
	}

	// Override with environment variables if set
	if timeout := os.Getenv("ANTHROPIC_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = d
		}
	}

	return cfg
}

// logUsageMetrics demonstrates how to track and log usage metrics
func logUsageMetrics(session *ChatSession) {
	// In a real application, you would track these metrics
	// For example: Prometheus, Datadog, or custom logging

	fmt.Println("Metrics would be tracked here in a production application:")
	fmt.Println("- Total API calls")
	fmt.Println("- Token usage (input/output)")
	fmt.Println("- Response times")
	fmt.Println("- Error rates")
	fmt.Println("- Model usage by type")
}

// gracefulShutdown handles cleanup on application exit
func gracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	// Perform cleanup here
	// For example: close database connections, flush logs, etc.
	fmt.Println("\nShutting down gracefully...")

	// Wait for context to be cancelled or timeout
	<-ctx.Done()
	fmt.Println("Shutdown complete")
}
