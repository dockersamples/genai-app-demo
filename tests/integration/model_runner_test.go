package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestModelRunnerIntegration tests connectivity to the Docker Model Runner
// and basic operations like listing, creating, and removing models
func TestModelRunnerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Docker Model Runner test in short mode")
	}

	// Start a Socat container to forward traffic to model-runner.docker.internal
	ctx := context.Background()
	socatContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "alpine/socat",
			Cmd:   []string{"tcp-listen:8080,fork,reuseaddr", "tcp:model-runner.docker.internal:80"},
			ExposedPorts: []string{
				"8080/tcp",
			},
			WaitingFor: wait.ForListeningPort("8080/tcp"),
		},
		Started: true,
	})

	if err != nil {
		t.Fatalf("Failed to start socat container: %s", err)
	}

	defer func() {
		if err := socatContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %s", err)
		}
	}()

	// Get the mapped port and host
	mappedPort, err := socatContainer.MappedPort(ctx, "8080")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	host, err := socatContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get host: %s", err)
	}

	// Create the base URL for the API
	baseURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
	client := &http.Client{Timeout: 10 * time.Second}

	// Test 1: GET /models endpoint
	t.Log("Testing GET /models endpoint...")
	resp, err := client.Get(baseURL + "/models")
	if err != nil {
		t.Fatalf("Failed to call /models: %s", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for /models")
	
	// Parse initial models list
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}
	
	var initialModels []map[string]interface{}
	if err := json.Unmarshal(body, &initialModels); err != nil {
		t.Fatalf("Failed to parse JSON response: %s", err)
	}
	
	// Test 2: POST /models/create endpoint to create a model
	t.Log("Testing POST /models/create endpoint...")
	modelName := "ignaciolopezluna020/llama3.2:1b"
	createModelReq, err := http.NewRequest(
		"POST", 
		baseURL+"/models/create", 
		strings.NewReader(fmt.Sprintf(`{"from": "%s"}`, modelName)),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}
	createModelReq.Header.Set("Content-Type", "application/json")
	
	createResp, err := client.Do(createModelReq)
	if err != nil {
		t.Fatalf("Failed to call /models/create: %s", err)
	}
	defer createResp.Body.Close()
	
	// Check if model creation was successful
	// Note: This might take a while if the model needs to be downloaded
	assert.Equal(t, http.StatusOK, createResp.StatusCode, "Expected status code 200 for /models/create")
	
	// Wait a bit for the model to be fully loaded
	time.Sleep(2 * time.Second)
	
	// Test 3: GET /models again to verify the model was added
	t.Log("Testing GET /models to verify model was added...")
	verifyResp, err := client.Get(baseURL + "/models")
	if err != nil {
		t.Fatalf("Failed to call /models for verification: %s", err)
	}
	defer verifyResp.Body.Close()
	
	assert.Equal(t, http.StatusOK, verifyResp.StatusCode, "Expected status code 200 for /models verification")
	
	verifyBody, err := io.ReadAll(verifyResp.Body)
	if err != nil {
		t.Fatalf("Failed to read verification response body: %s", err)
	}
	
	var modelsAfterCreation []map[string]interface{}
	if err := json.Unmarshal(verifyBody, &modelsAfterCreation); err != nil {
		t.Fatalf("Failed to parse JSON verification response: %s", err)
	}
	
	// Check that the models list increased by 1
	expectedModelsCount := len(initialModels) + 1
	assert.Equal(t, expectedModelsCount, len(modelsAfterCreation), "Expected models count to increase by 1")
	
	// Test 4: DELETE /models/{model_name} to clean up
	t.Log("Testing DELETE /models endpoint...")
	deleteReq, err := http.NewRequest("DELETE", baseURL+"/models/"+modelName, nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %s", err)
	}
	
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("Failed to call DELETE /models/%s: %s", modelName, err)
	}
	defer deleteResp.Body.Close()
	
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode, "Expected status code 200 for DELETE /models")
	
	// Test 5: Final verification that the model was removed
	t.Log("Final verification that model was removed...")
	finalResp, err := client.Get(baseURL + "/models")
	if err != nil {
		t.Fatalf("Failed to call /models for final verification: %s", err)
	}
	defer finalResp.Body.Close()
	
	finalBody, err := io.ReadAll(finalResp.Body)
	if err != nil {
		t.Fatalf("Failed to read final verification response body: %s", err)
	}
	
	var finalModels []map[string]interface{}
	if err := json.Unmarshal(finalBody, &finalModels); err != nil {
		t.Fatalf("Failed to parse JSON final verification response: %s", err)
	}
	
	assert.Equal(t, len(initialModels), len(finalModels), "Expected final models count to match initial count")
}