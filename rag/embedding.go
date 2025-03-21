package rag

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// EmbeddingClient is a client for generating embeddings
type EmbeddingClient struct {
	client *openai.Client
}

// NewEmbeddingClient creates a new embedding client
func NewEmbeddingClient() *EmbeddingClient {
	baseURL := os.Getenv("BASE_URL")
	apiKey := os.Getenv("API_KEY")

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)

	return &EmbeddingClient{client: client}
}

// GenerateEmbedding generates an embedding for the given text
func (c *EmbeddingClient) GenerateEmbedding(ctx context.Context, text string, model string) ([]float32, error) {
	// If no model specified, use a default
	if model == "" {
		model = "text-embedding-ada-002"
	}

	// This is a simplified version - in a real implementation,
	// you would call the actual OpenAI embedding API
	// For now, we'll simulate an embedding

	// In a real implementation, you would use code like this:
	/*
	   params := openai.EmbeddingCreationParams{
	       Model: openai.F(model),
	       Input: openai.F([]string{text}),
	   }
	   response, err := c.client.Embeddings.Create(ctx, params)
	   if err != nil {
	       return nil, fmt.Errorf("failed to generate embedding: %w", err)
	   }
	   return response.Data[0].Embedding, nil
	*/

	// For now, return a mock embedding
	return []float32{0.1, 0.2, 0.3, 0.4, 0.5}, nil
}

// GenerateDocumentEmbedding generates an embedding for a document
func (c *EmbeddingClient) GenerateDocumentEmbedding(ctx context.Context, doc Document, model string) ([]float32, error) {
	// Combine title and content for embedding
	text := fmt.Sprintf("%s\n%s", doc.Title, doc.Content)
	return c.GenerateEmbedding(ctx, text, model)
}

// GenerateQueryEmbedding generates an embedding for a query
func (c *EmbeddingClient) GenerateQueryEmbedding(ctx context.Context, query string, model string) ([]float32, error) {
	return c.GenerateEmbedding(ctx, query, model)
}
