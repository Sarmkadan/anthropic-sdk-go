package benchmarks

import (
	"context"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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
				Model:     anthropic.ModelClaudeOpus4_6,
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
				System: []anthropic.TextBlockParam{
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
				Model:     anthropic.ModelClaudeOpus4_6,
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
			_, _ = client.Models.List(context.Background(), anthropic.ModelListParams{})
		}
	})

	b.Run("ModelRetrieve", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Models.Get(context.Background(), "claude-opus-4-6", anthropic.ModelGetParams{})
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
				Model:  anthropic.ModelClaudeOpus4_6,
				Prompt: "What is a quaternion?",
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
			_ = anthropic.TextBlockParam{
				Text: "Test content",
			}
		}
	})

	b.Run("ImageBlockWithBase64", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.ImageBlockParam{
				Source: anthropic.ImageBlockParamSourceUnion{
					OfBase64: &anthropic.Base64ImageSourceParam{
						Data:      "fake-image-data",
						MediaType: "image/jpeg",
					},
				},
			}
		}
	})

	b.Run("DocumentBlockWithBase64", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.DocumentBlockParam{
				Source: anthropic.DocumentBlockParamSourceUnion{
					OfBase64: &anthropic.Base64PDFSourceParam{
						Data:      "fake-pdf-data",
						MediaType: "application/pdf",
					},
				},
				Title: "Test Document",
			}
		}
	})
}

// BenchmarkMessageValidation measures the overhead of message validation and parameter processing
func BenchmarkMessageValidation(b *testing.B) {
	params := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_6,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Test message"),
		),
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = params
	}
}

// BenchmarkMemory measures memory allocations for critical operations
func BenchmarkMemory(b *testing.B) {
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
				Model:     anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Test message content"),
					),
				},
			}
		}
	})

	b.Run("BlockTypesMemory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.TextBlockParam{
				Text: "Test content",
			}
		}
	})
}

// BenchmarkThroughput measures throughput of message operations
func BenchmarkThroughput(b *testing.B) {
	b.Run("MessageParamsPerSecond", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageNewParams{
				Model:     anthropic.ModelClaudeOpus4_6,
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
					Model:     anthropic.ModelClaudeOpus4_6,
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
				Model:     anthropic.ModelClaudeOpus4_6,
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
				Model:     anthropic.ModelClaudeSonnet4_5,
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

// BenchmarkBlockSizes measures performance with different block sizes
func BenchmarkBlockSizes(b *testing.B) {
	b.Run("SmallTextBlock", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewTextBlock("Small")
		}
	})

	b.Run("MediumTextBlock", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewTextBlock("This is a medium-sized text block with moderate content.")
		}
	})

	b.Run("LargeTextBlock", func(b *testing.B) {
		b.ResetTimer()
		largeText := "This is a very large text block with significant content to test performance with larger data sizes. "
		largeText += "It contains multiple sentences and should provide a good measure of how the library handles larger payloads. "
		largeText += "Performance with large blocks is important for applications that process documents or large pieces of text."
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewTextBlock(largeText)
		}
	})
}

// BenchmarkMessageBatch measures performance of message batch operations
func BenchmarkMessageBatch(b *testing.B) {
	b.Run("CreateMessageBatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.MessageBatchNewParams{
				Model:     anthropic.ModelClaudeOpus4_6,
				MaxTokens: 1024,
				BatchSize: 10,
				Messages: []anthropic.MessageParam{
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Message 1"),
					),
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Message 2"),
					),
					anthropic.NewUserMessage(
						anthropic.NewTextBlock("Message 3"),
					),
				},
			}
		}
	})
}

// BenchmarkClientConfiguration measures performance of different client configurations
func BenchmarkClientConfiguration(b *testing.B) {
	b.Run("ClientWithTimeout", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient(
				option.WithAPIKey("test-api-key-123456789"),
				option.WithHTTPClientTimeout(30),
			)
		}
	})

	b.Run("ClientWithBaseURL", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = anthropic.NewClient(
				option.WithAPIKey("test-api-key-123456789"),
				option.WithBaseURL("https://api.anthropic.com"),
			)
		}
	})
}
