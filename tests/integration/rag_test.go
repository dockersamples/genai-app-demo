package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/neo4j"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Document represents a document to add to the knowledge base
type Document struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url,omitempty"`
}

// DocumentResponse is the response from adding a document
type DocumentResponse struct {
	ID string `json:"id"`
}

// TestNeo4jRagIntegration tests the integration between the API and Neo4j
func TestNeo4jRagIntegration(t *testing.T) {
	// Skip if CI environment doesn't have docker
	if os.Getenv("CI") != "" && os.Getenv("DOCKER_AVAILABLE") == "" {
		t.Skip("Skipping test in CI environment without Docker")
	}

	// Start Neo4j container
	ctx := context.Background()
	
	neo4jContainer, err := neo4j.RunContainer(ctx,
		testcontainers.WithImage("neo4j:5.15.0-community"),
		neo4j.WithAdminPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Started.")
		),
	)
	if err != nil {
		t.Fatalf("Failed to start Neo4j container: %v", err)
	}

	// Make sure to clean up the container
	defer func() {
		if err := neo4jContainer.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate Neo4j container: %v", err)
		}
	}()

	// Get Neo4j connection details
	neo4jURI, err := neo4jContainer.URI(ctx)
	if err != nil {
		t.Fatalf("Failed to get Neo4j URI: %v", err)
	}

	// Start the backend container
	backendContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "genai-app-demo-backend:test",
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"BASE_URL":       "http://localhost:8080/v1",
				"MODEL":         "gpt-3.5-turbo",
				"API_KEY":       "test-key",
				"NEO4J_URI":     neo4jURI,
				"NEO4J_USERNAME": "neo4j",
				"NEO4J_PASSWORD": "password",
				"RAG_ENABLED":   "true",
			},
			WaitingFor: wait.ForHTTP("/health").WithPort("8080/tcp"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Failed to start backend container: %v", err)
	}

	// Clean up the backend container
	defer func() {
		if err := backendContainer.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate backend container: %v", err)
		}
	}()

	// Get the backend URL
	backendHost, err := backendContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get backend host: %v", err)
	}

	backendPort, err := backendContainer.MappedPort(ctx, "8080/tcp")
	if err != nil {
		t.Fatalf("Failed to get backend port: %v", err)
	}

	backendURL := fmt.Sprintf("http://%s:%s", backendHost, backendPort.Port())

	// Test adding a document
	doc := Document{
		Title:   "Test Document",
		Content: "This is a test document about artificial intelligence. AI is transforming our world.",
		URL:     "https://example.com/test",
	}

	// Add the document to the knowledge base
	docID, err := addDocument(ctx, backendURL, doc)
	require.NoError(t, err, "Should add document without error")
	assert.NotEmpty(t, docID, "Document ID should not be empty")

	// Verify the document was added to Neo4j
	assert.Eventually(t, func() bool {
		return verifyDocumentInNeo4j(ctx, neo4jURI, docID)
	}, 10*time.Second, 1*time.Second, "Document should be added to Neo4j")

	// Test the RAG functionality with a query related to the document
	t.Log("Testing RAG functionality with related query")
	response, err := sendChatRequest(ctx, backendURL, "Tell me about AI")
	require.NoError(t, err, "Should send chat request without error")
	// In a real test, we would verify the response contains information from the document
	assert.NotEmpty(t, response, "Response should not be empty")
}

// Helper function to add a document
func addDocument(ctx context.Context, baseURL string, doc Document) (string, error) {
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal document: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/documents",
		bytes.NewBuffer(docBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response DocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.ID, nil
}

// Helper function to verify a document exists in Neo4j
func verifyDocumentInNeo4j(ctx context.Context, uri string, docID string) bool {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		return false
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "MATCH (d:Document {id: $id}) RETURN count(d) as count", map[string]interface{}{
		"id": docID,
	})
	if err != nil {
		return false
	}

	record, err := result.Single(ctx)
	if err != nil {
		return false
	}

	count, _ := record.Get("count")
	return count.(int64) > 0
}

// Helper function to send a chat request
func sendChatRequest(ctx context.Context, baseURL string, message string) (string, error) {
	reqBody := map[string]interface{}{
		"messages": []map[string]string{},
		"message":  message,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/chat",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// For event stream responses, we would need to parse the streaming response
	// For testing purposes, we'll just read the entire response
	responseBytes := make([]byte, 4096) // Read up to 4KB
	n, err := resp.Body.Read(responseBytes)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(responseBytes[:n]), nil
}
