#!/bin/bash
# Setup script for Neo4j testcontainer tests

# Ensure we have the right dependencies
echo "Installing Neo4j Go Driver..."
go get github.com/neo4j/neo4j-go-driver/v5/neo4j@latest

# Update dependencies
echo "Updating dependencies..."
go mod download && go mod tidy

# Navigate to tests directory
cd tests/
go mod download && go mod tidy

# Run the Neo4j vector tests
echo "Running Neo4j Vector Tests..."
go test -v ./integration -run TestNeo4jVectorDatabase

echo "Setup complete!"
