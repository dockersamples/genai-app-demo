package rag

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
)

// RAGManager handles the retrieval-augmented generation workflow
type RAGManager struct {
	neo4jClient     *Neo4jClient
	embeddingClient *EmbeddingClient
	enabled         bool
	contextLimit    int
}

// NewRAGManager creates a new RAG manager
func NewRAGManager() (*RAGManager, error) {
	// Check if RAG is enabled
	enabled := true
	ragEnabledStr := os.Getenv("RAG_ENABLED")
	if ragEnabledStr != "" {
		var err error
		enabled, err = strconv.ParseBool(ragEnabledStr)
		if err != nil {
			log.Printf("Warning: Invalid RAG_ENABLED value, defaulting to true: %v", err)
			enabled = true
		}
	}

	if !enabled {
		return &RAGManager{enabled: false}, nil
	}

	// Parse context limit
	contextLimit := 5
	contextLimitStr := os.Getenv("RAG_CONTEXT_LIMIT")
	if contextLimitStr != "" {
		var err error
		contextLimit, err = strconv.Atoi(contextLimitStr)
		if err != nil {
			log.Printf("Warning: Invalid RAG_CONTEXT_LIMIT value, defaulting to 5: %v", err)
			contextLimit = 5
		}
	}

	// Create Neo4j client
	neo4jClient, err := NewNeo4jClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j client: %w", err)
	}

	// Create embedding client
	embeddingClient := NewEmbeddingClient()

	return &RAGManager{
		neo4jClient:     neo4jClient,
		embeddingClient: embeddingClient,
		enabled:         enabled,
		contextLimit:    contextLimit,
	}, nil
}

// Close closes all resources used by the RAG manager
func (m *RAGManager) Close(ctx context.Context) error {
	if !m.enabled || m.neo4jClient == nil {
		return nil
	}
	return m.neo4jClient.Close(ctx)
}

// IsEnabled returns whether RAG is enabled
func (m *RAGManager) IsEnabled() bool {
	return m.enabled
}

// EnhancePromptWithContext enhances a prompt with relevant context from the knowledge base
func (m *RAGManager) EnhancePromptWithContext(ctx context.Context, query string) (string, error) {
	if !m.enabled {
		return query, nil
	}

	// Get relevant documents
	documents, err := m.neo4jClient.GetRelevantDocuments(ctx, query, m.contextLimit)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve relevant documents: %w", err)
	}

	if len(documents) == 0 {
		// No relevant documents found, return the original query
		return query, nil
	}

	// Format the context for the prompt
	context := FormatContextForPrompt(documents)
	
	// Combine the context with the query
	prompt := fmt.Sprintf(
		"I'll provide you with some relevant context to help answer the following question.\n\nQuestion: %s\n\n%s\n\nPlease provide an answer based on the context provided. If the context doesn't contain relevant information, say so and try to provide a general answer.",
		query,
		context,
	)

	return prompt, nil
}

// AddDocument adds a document to the knowledge base
func (m *RAGManager) AddDocument(ctx context.Context, doc Document) error {
	if !m.enabled {
		return fmt.Errorf("RAG is not enabled")
	}

	// Generate embedding for the document
	_, err := m.embeddingClient.GenerateDocumentEmbedding(ctx, doc, "")
	if err != nil {
		return fmt.Errorf("failed to generate document embedding: %w", err)
	}

	// Add the document to Neo4j
	err = m.neo4jClient.AddDocument(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to add document to Neo4j: %w", err)
	}

	return nil
}
