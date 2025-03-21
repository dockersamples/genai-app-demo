package rag

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jClient is a wrapper around Neo4j Driver for RAG operations
type Neo4jClient struct {
	driver neo4j.DriverWithContext
}

// NewNeo4jClient creates a new Neo4j client for RAG operations
func NewNeo4jClient() (*Neo4jClient, error) {
	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	if uri == "" || username == "" || password == "" {
		return nil, fmt.Errorf("NEO4J_URI, NEO4J_USERNAME, and NEO4J_PASSWORD must be set")
	}

	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	return &Neo4jClient{driver: driver}, nil
}

// Close closes the Neo4j driver
func (c *Neo4jClient) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

// Document represents a document in the knowledge graph
type Document struct {
	ID          string
	Title       string
	Content     string
	URL         string
	EmbeddingID string
}

// RetrieverResult represents a result from the retriever
type RetrieverResult struct {
	Document Document
	Score    float64
}

// GetRelevantDocuments retrieves documents relevant to the query from Neo4j
func (c *Neo4jClient) GetRelevantDocuments(ctx context.Context, query string, limit int) ([]RetrieverResult, error) {
	if limit <= 0 {
		limit = 5 // Default limit
	}

	// For now, we'll do a simple keyword-based search
	// In a real implementation, we'd use vector similarity search with embeddings
	keywords := strings.Split(query, " ")
	cypher := `
		MATCH (d:Document)
		WHERE ANY(keyword IN $keywords WHERE d.content CONTAINS keyword)
		RETURN d.id as id, d.title as title, d.content as content, d.url as url, d.embeddingId as embeddingId
		ORDER BY size([keyword IN $keywords WHERE d.content CONTAINS keyword]) DESC
		LIMIT $limit
	`

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, cypher, map[string]interface{}{
		"keywords": keywords,
		"limit":    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var documents []RetrieverResult
	for result.Next(ctx) {
		record := result.Record()
		
		id, _ := record.Get("id")
		title, _ := record.Get("title")
		content, _ := record.Get("content")
		url, _ := record.Get("url")
		embeddingId, _ := record.Get("embeddingId")

		doc := Document{
			ID:          toString(id),
			Title:       toString(title),
			Content:     toString(content),
			URL:         toString(url),
			EmbeddingID: toString(embeddingId),
		}

		// Calculate a simple score based on keyword matches
		// In a real implementation, this would be based on vector similarity
		matchCount := 0
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(doc.Content), strings.ToLower(keyword)) {
				matchCount++
			}
		}
		score := float64(matchCount) / float64(len(keywords))

		documents = append(documents, RetrieverResult{Document: doc, Score: score})
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error during result iteration: %w", err)
	}

	return documents, nil
}

// FormatContextForPrompt formats the retrieved documents into a context string for the prompt
func FormatContextForPrompt(results []RetrieverResult) string {
	var sb strings.Builder
	sb.WriteString("\nRelevant context:\n")

	for i, result := range results {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, result.Document.Title))
		
		// Truncate content if too long
		content := result.Document.Content
		if len(content) > 500 {
			content = content[:497] + "..."
		}
		
		sb.WriteString(content)
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// AddDocument adds a document to the knowledge graph
func (c *Neo4jClient) AddDocument(ctx context.Context, doc Document) error {
	cypher := `
		MERGE (d:Document {id: $id})
		SET d.title = $title,
		    d.content = $content,
		    d.url = $url,
		    d.embeddingId = $embeddingId
		RETURN d
	`

	session := c.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx, cypher, map[string]interface{}{
		"id":          doc.ID,
		"title":       doc.Title,
		"content":     doc.Content,
		"url":         doc.URL,
		"embeddingId": doc.EmbeddingID,
	})

	if err != nil {
		return fmt.Errorf("failed to add document: %w", err)
	}

	return nil
}

// Helper function to convert interface{} to string
func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
