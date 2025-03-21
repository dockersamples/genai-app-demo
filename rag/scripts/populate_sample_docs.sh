#!/bin/bash

# This script adds sample documents to the Neo4j knowledge base

BACKEND_URL=${BACKEND_URL:-http://localhost:8080}

echo "Adding sample documents to knowledge base at $BACKEND_URL"

# Add document 1
echo "Adding document: What is Generative AI?"
curl -X POST "$BACKEND_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "What is Generative AI?",
    "content": "Generative AI refers to artificial intelligence systems that can generate new content, including text, images, audio, code, and more. These systems are trained on large datasets and learn to create original outputs that resemble the training data. Unlike traditional AI systems that analyze or classify existing data, generative AI creates new content. Popular examples include large language models (LLMs) like GPT-4, image generators like DALL-E and Stable Diffusion, and music creators like MusicLM. These technologies have applications in creative work, content production, software development, customer service, and many other fields.",
    "url": "https://example.com/generative-ai"
  }'
echo ""

# Add document 2
echo "Adding document: Retrieval-Augmented Generation (RAG)"
curl -X POST "$BACKEND_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Retrieval-Augmented Generation (RAG)",
    "content": "Retrieval-Augmented Generation (RAG) is an AI framework that enhances large language models by incorporating external knowledge retrieval. Instead of relying solely on the knowledge encoded in the model\'s parameters during training, RAG systems retrieve relevant information from external knowledge sources when generating responses. This approach combines the strengths of retrieval-based and generation-based systems. The benefits include access to more up-to-date information, reduced hallucinations, improved factual accuracy, and the ability to cite sources. RAG typically involves a knowledge base, a retrieval component, and a generation component working together to produce more reliable and informative responses.",
    "url": "https://example.com/rag-explained"
  }'
echo ""

# Add document 3
echo "Adding document: Graph Databases and Neo4j"
curl -X POST "$BACKEND_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Graph Databases and Neo4j",
    "content": "Graph databases are specialized NoSQL databases designed to store and navigate relationships between data entities. Unlike traditional relational databases that use tables, graph databases use nodes, edges, and properties to represent and store data. Neo4j is one of the most popular graph database systems, known for its high performance and scalability. It uses the property graph model, where nodes and relationships can both contain properties. Neo4j supports the Cypher query language, which is specifically designed for working with graph data. Graph databases excel at use cases like social networks, recommendation engines, fraud detection, and knowledge graphs, where relationship patterns and connections between data points are as important as the data itself.",
    "url": "https://example.com/neo4j-intro"
  }'
echo ""

# Add document 4
echo "Adding document: Knowledge Graphs in AI Systems"
curl -X POST "$BACKEND_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Knowledge Graphs in AI Systems",
    "content": "Knowledge graphs represent information as a network of entities and their relationships, encoded with semantic metadata. In AI systems, they serve as structured knowledge repositories that can be queried and navigated. Unlike traditional vector databases that store embeddings for similarity search, knowledge graphs capture explicit relationships between concepts, enabling more complex reasoning. When integrated with generative AI, knowledge graphs can provide contextual information, support multi-hop reasoning, and help explain AI decisions through transparent fact paths. They excel at representing hierarchical information, domain expertise, and complex relationships between entities, making them powerful tools for enhancing AI with structured knowledge.",
    "url": "https://example.com/knowledge-graphs-ai"
  }'
echo ""

# Add document 5
echo "Adding document: Vector Embeddings for Semantic Search"
curl -X POST "$BACKEND_URL/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vector Embeddings for Semantic Search",
    "content": "Vector embeddings transform text, images, or other data into numerical vectors in a high-dimensional space, where similar items are positioned closer together. In semantic search applications, these embeddings capture the meaning and context of content rather than just keywords. When a user submits a query, it is converted into an embedding vector and compared with document embeddings using similarity metrics like cosine similarity. This approach enables search systems to find conceptually related content even when exact keywords don\'t match. Modern embedding models like those from OpenAI, Cohere, and open-source alternatives have dramatically improved the quality of semantic search, allowing for more intuitive and contextually aware information retrieval.",
    "url": "https://example.com/vector-embeddings"
  }'
echo ""

echo "Finished adding sample documents to knowledge base"
