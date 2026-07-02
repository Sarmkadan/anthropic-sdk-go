package benchmarks

import (
	"context"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/benchmarkdotnet/diagnoser"
)

// BenchmarkClientCreation measures the time to create new client instances
// This tests the overhead of client initialization including credential resolution
func BenchmarkClientCreation(b *testing.B) {
	b.Run("DefaultClient", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient()
		}
	})

	b.Run("ClientWithAPIKey", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient(
				option.WithAPIKey("test-api-key-123456789"),
			)
		}
	})

	b.Run("ClientWithConfig", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient(
				option.WithAPIKey("test-api-key-123456789"),
				option.WithEnvironmentProduction(),
			)
		}
	})
}

// BenchmarkMessageCreation measures the time to create message requests
// This tests the overhead of request parameter construction and validation
func BenchmarkMessageCreation(b *testing.B) {
	b.Run("SimpleMessage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("What is a quaternion?"),
					),
				},
			}
		}
	})

	b.Run("MessageWithSystemPrompt", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model: anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				System: []anthropic.SystemPromptBlock{
					anthropic.NewTextBlock("You are a helpful assistant that explains mathematical concepts."),
				},
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("What is a quaternion?"),
					),
				},
			}
		}
	})

	b.Run("MessageWithMultipleBlocks", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Explain the following concepts: "),
						anthropic.NewTextBlock("- Vectors"),
						anthropic.NewTextBlock("- Matrices"),
						anthropic.NewTextBlock("- Tensors"),
					),
				},
			}
		}
	})
}

// BenchmarkMessageServiceCreation measures the time to create message service instances
func BenchmarkMessageServiceCreation(b *testing.B) {
	b.Run("NewMessageService", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewMessageService(
				option.WithAPIKey("test-api-key-123456789"),
			)
		}
	})
}

// BenchmarkModelOperations measures performance of model-related operations
func BenchmarkModelOperations(b *testing.B) {
	client := anthropic.NewClient(
		option.WithAPIKey("test-api-key-123456789"),
	)

	b.Run("ModelList", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Models.List(context.Background())
		}
	})

	b.Run("ModelRetrieve", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Models.Get(context.Background(), "claude-opus-4-6")
		}
	})
}

// BenchmarkCompletionOperations measures performance of completion operations
func BenchmarkCompletionOperations(b *testing.B) {
	client := anthropic.NewClient(
		option.WithAPIKey("test-api-key-123456789"),
	)

	b.Run("CompletionCreate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Completions.New(context.Background(), anthropic.CompletionNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Prompt:   "What is a quaternion?",
			})
		}
	})
}

// BenchmarkMessageParamConstruction measures the overhead of creating different message parameter types
func BenchmarkMessageParamConstruction(b *testing.B) {
	b.Run("UserMessageWithText", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewUserMessage(
				anthropic.NewTextBlock("What is a quaternion?"),
			)
		}
	})

	b.Run("AssistantMessageWithText", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewAssistantMessage(
				anthropic.NewTextBlock("A quaternion is a mathematical concept that extends complex numbers."),
			)
		}
	})

	b.Run("ToolMessage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewToolMessage(
				"tool-123",
				anthropic.NewTextBlock("42"),
			)
		}
	})

	b.Run("TextBlockCreation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewTextBlock("This is a test message block with some content to measure performance.")
		}
	})
}

// BenchmarkBlockTypes measures performance of different block types
func BenchmarkBlockTypes(b *testing.B) {
	b.Run("TextBlock", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.TextBlock{Type: "text", Text: "Test content"}
		}
	})

	b.Run("ImageBlock", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.ImageBlock{
				Type:      "image",
				MediaType: "image/jpeg",
				Data:      []byte("fake-image-data"),
			}
		}
	})

	b.Run("DocumentBlock", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.DocumentBlock{
				Type:      "document",
				MediaType: "application/pdf",
				Data:      []byte("fake-pdf-data"),
			}
		}
	})
}

// BenchmarkMessageValidation measures the overhead of message validation and parameter processing
func BenchmarkMessageValidation(b *testing.B) {
	params := anthropic.MessageNewParams{
		Model:    anthropic.ModelClaudeOpus4_6,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Test message"),
		),
	},
}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = params.Validate()
	}
}

// BenchmarkMemoryDiagnostics enables memory allocation tracking for critical operations
// This benchmark will show memory allocations for critical operations
func BenchmarkMemoryDiagnostics(b *testing.B) {
	b.Run("ClientCreationMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient(
				option.WithAPIKey("test-api-key-123456789"),
			)
		}
	})

	b.Run("MessageParamsMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test message content"),
				),
			},
		}
	})

	b.Run("BlockTypesMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.TextBlock{Type: "text", Text: "Test content"}
		}
	})
}

// BenchmarkThroughput measures throughput of message operations
func BenchmarkThroughput(b *testing.B) {
	b.Run("MessageParamsPerSecond", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test message"),
					),
				},
			}
		}
	})

	b.Run("MultipleMessagesPerSecond", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			params := make([]anthropic.MessageParam, 100)
			for j := 0; j < 100; j++ {
				params[j] = anthropic.NewUserMessage(
					anthropic.NewTextBlock("Test message "),
				)
			}
			_ = params
		}
	})
}

// BenchmarkConcurrentOperations measures performance under concurrent load
func BenchmarkConcurrentOperations(b *testing.B) {
	b.Run("ConcurrentClientCreation", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = anthropic.NewClient(
					option.WithAPIKey("test-api-key-123456789"),
				)
			}
		})
	})

	b.Run("ConcurrentMessageCreation", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = anthropic.MessageNewParams{
					Model:    anthropic.ModelClaudeOpus4_6,
					MaxTokens: 1024,
					Messages: []anthropic.MessageParam{
						anthropic.NewUserMessage(
							anthropic.NewTextBlock("Test message"),
						),
					},
				}
			}
		})
	})
}

// BenchmarkDifferentModels measures performance across different model types
func BenchmarkDifferentModels(b *testing.B) {
	b.Run("ModelClaudeOpus4_6", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test"),
					),
				},
			}
		}
	})

	b.Run("ModelClaudeSonnet4_5", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:    anthropic.ModelClaudeSonnet4_5,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test"),
					),
				},
			}
		}
	})
}

// BenchmarkErrorHandling measures performance of error scenarios
func BenchmarkErrorHandling(b *testing.B) {
	b.Run("InvalidMessageParams", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			params := anthropic.MessageNewParams{
				Model:    "", // Invalid empty model
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test"),
					),
				},
			}
			_ = params.Validate()
		}
	})
}