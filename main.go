package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/username/project/rag" // Import our RAG package
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []Message `json:"messages"`
	Message  string    `json:"message"`
}

// Document request for adding documents to the knowledge base
type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url,omitempty"`
}

func main() {
	baseURL := os.Getenv("BASE_URL")
	model := os.Getenv("MODEL")
	apiKey := os.Getenv("API_KEY")

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)

	// Initialize the RAG manager
	ragManager, err := rag.NewRAGManager()
	if err != nil {
		log.Fatalf("Failed to initialize RAG manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := ragManager.Close(ctx); err != nil {
			log.Printf("Error closing RAG manager: %v", err)
		}
	}()

	log.Printf("RAG enabled: %v", ragManager.IsEnabled())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Endpoint for adding documents to the knowledge base
	http.HandleFunc("/documents", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req DocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Create a document and add it to the knowledge base
		doc := rag.Document{
			ID:      uuid.New().String(),
			Title:   req.Title,
			Content: req.Content,
			URL:     req.URL,
		}

		ctx := r.Context()
		err := ragManager.AddDocument(ctx, doc)
		if err != nil {
			log.Printf("Error adding document: %v", err)
			http.Error(w, "Failed to add document", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": doc.ID})
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ctx := r.Context()

		var messages []openai.ChatCompletionMessageParamUnion
		for _, msg := range req.Messages {
			var message openai.ChatCompletionMessageParamUnion
			switch msg.Role {
			case "user":
				message = openai.UserMessage(msg.Content)
			case "assistant":
				message = openai.AssistantMessage(msg.Content)
			}

			messages = append(messages, message)
		}

		// Process the user's message with RAG if enabled
		userMsg := req.Message
		if ragManager.IsEnabled() {
			enhancedMsg, err := ragManager.EnhancePromptWithContext(ctx, userMsg)
			if err != nil {
				log.Printf("Warning: Failed to enhance prompt with RAG: %v", err)
				// Fall back to the original message
			} else {
				userMsg = enhancedMsg
			}
		}

		param := openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Model:    openai.F(model),
		}

		// Adds the user message to the conversation
		param.Messages.Value = append(param.Messages.Value, openai.UserMessage(userMsg))
		stream := client.Chat.Completions.NewStreaming(ctx, param)

		for stream.Next() {
			chunk := stream.Current()

			// Stream each chunk as it arrives
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				_, err := fmt.Fprintf(w, "%s", chunk.Choices[0].Delta.Content)
				if err != nil {
					fmt.Printf("Error writing to stream: %v\n", err)
					return
				}
				w.(http.Flusher).Flush()
			}
		}

		if err := stream.Err(); err != nil {
			fmt.Printf("Error in stream: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
