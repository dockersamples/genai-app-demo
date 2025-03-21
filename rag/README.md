# Neo4j RAG Implementation

This directory contains the implementation of Retrieval-Augmented Generation (RAG) using Neo4j as the knowledge graph database.

## Overview

This RAG implementation enhances the GenAI chat application by grounding responses in a knowledge base stored in Neo4j. The system uses a hybrid retrieval approach that leverages both keyword-based search and vector similarity (via embeddings) to find the most relevant documents for a user query.

## Components

### 1. Neo4j Client (`neo4j.go`)

Handles interactions with the Neo4j database:

- Connection management with Neo4j
- Document storage and retrieval
- Document search functionality

### 2. Embedding Client (`embedding.go`)

Generates vector embeddings for documents and queries:

- Text embedding generation using the OpenAI API
- Document embedding functionality
- Query embedding functionality

### 3. RAG Manager (`manager.go`)

Coordinates the RAG workflow:

- Configures RAG behavior based on environment variables
- Manages document ingestion
- Processes user queries to retrieve relevant context
- Enhances prompts with retrieved context

## Data Model

The Neo4j data model consists of the following structure:

```cypher
(Document)-[:HAS_EMBEDDING]->(Embedding)
```

Document nodes contain:
- `id`: Unique identifier for the document
- `title`: Document title
- `content`: Document text content
- `url`: Optional URL source

## Integration Flow

1. **Document Ingestion**:
   - Document is received via API endpoint
   - Vector embedding is generated
   - Document and embedding are stored in Neo4j

2. **Query Processing**:
   - User query is analyzed
   - Relevant documents are retrieved from Neo4j
   - LLM prompt is enhanced with retrieved context
   - Enhanced prompt is sent to LLM

## Configuration

The RAG system can be configured through environment variables:

- `NEO4J_URI`: Connection URI for Neo4j (default: `neo4j://neo4j:7687`)
- `NEO4J_USERNAME`: Neo4j username (default: `neo4j`)
- `NEO4J_PASSWORD`: Neo4j password (default: `password`)
- `RAG_ENABLED`: Enable/disable RAG functionality (default: `true`)
- `RAG_CONTEXT_LIMIT`: Maximum number of documents to retrieve (default: `5`)

## API Endpoints

### Add Document

```
POST /documents
```

Request body:
```json
{
  "title": "Document Title",
  "content": "Document content text...",
  "url": "https://example.com/source" (optional)
}
```

Response:
```json
{
  "id": "generated-document-id"
}
```

## Testing

The RAG implementation can be tested using the Testcontainers-based integration tests:

```bash
cd tests
go test -v ./integration -run TestNeo4jRagIntegration
```

## Future Improvements

1. **Advanced Vector Search**: Implement more sophisticated vector search algorithms
2. **Chunking Strategy**: Add document chunking for more granular context retrieval
3. **Multi-hop Reasoning**: Add graph traversal for answering complex questions requiring multiple documents
4. **Relevance Ranking**: Improve document ranking with better scoring algorithms
5. **Caching**: Add response caching for frequently asked questions
