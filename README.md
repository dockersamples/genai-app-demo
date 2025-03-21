# GenAI Application Demo

A modern chat application demonstrating integration of frontend technologies with local Large Language Models (LLMs).

## Overview

This project is a full-stack GenAI chat application that showcases how to build a Generative AI interface with a React frontend and Go backend, connected to Llama 3.2 running on Ollama.

## Architecture

The application consists of three main components:

1. **Frontend**: React TypeScript application providing a responsive chat interface
2. **Backend**: Go server that handles API requests and connects to the LLM
3. **Model Runner**: Ollama service running the Llama 3.2 (1B parameter) model

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │ >>> │   Backend   │ >>> │    Ollama   │
│  (React/TS) │     │    (Go)     │     │  (Llama 3.2)│
└─────────────┘     └─────────────┘     └─────────────┘
      :3000              :8080              :11434
```

## Features

- Real-time chat interface with message history
- Streaming AI responses (tokens appear as they're generated)
- Dockerized deployment for easy setup
- Local LLM integration (no cloud API dependencies)
- Cross-origin resource sharing (CORS) enabled
- Comprehensive integration tests using Testcontainers

## Prerequisites

- Docker and Docker Compose
- Git
- Go 1.19 or higher

## Quick Start

1. Clone this repository:
   ```bash
   git clone https://github.com/ajeetraina/genai-app-demo.git
   cd genai-app-demo
   ```

2. Start the application using Docker Compose:
   ```bash
   docker compose up -d -build
   ```

3. Access the frontend at [http://localhost:3000](http://localhost:3000)

## Development Setup

### Frontend

The frontend is a React TypeScript application using Vite:

```bash
cd frontend
npm install
npm run dev
```

### Backend

The Go backend can be run directly:

```bash
go mod download
go run main.go
```

Make sure to set the environment variables in `backend.env` or provide them directly.


## License

[MIT](https://github.com/dockersamples/genai-app-demo/blob/main/LICENSE)

## Contributing

Since this project is intended to support a specific use case guide, contributions are limited to bug fixes or security issues. If you have a question, feel free to open an issue!
