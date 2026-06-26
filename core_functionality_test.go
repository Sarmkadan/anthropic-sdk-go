// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestClientInitialization(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("test-api-key"),
	)
	// NewClient returns a struct, so we can't check for nil.
	// We'll check if options are populated instead.
	if len(client.Options) == 0 {
		t.Fatal("expected client options to be populated")
	}
}

func TestMessageParamsConstruction(t *testing.T) {
	params := anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "hello",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet4_5,
	}

	if params.MaxTokens != 1024 {
		t.Errorf("expected MaxTokens 1024, got %d", params.MaxTokens)
	}
	if len(params.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(params.Messages))
	}
}

func TestErrorHandling400(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       http.NoBody,
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 10,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestErrorHandling500(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       http.NoBody,
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 10,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModelEnumConstants(t *testing.T) {
	model := anthropic.ModelClaudeSonnet4_5
	if model == "" {
		t.Fatal("expected model constant to be set")
	}
}
