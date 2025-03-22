# GenAI Application Demo using Model Runner

A modern chat application demonstrating integration of frontend technologies with local Large Language Models (LLMs), with a focus on Docker Model Runner integration.

## Overview

This project is a full-stack GenAI chat application that showcases how to build a Generative AI interface with a React frontend and Go backend, connected to Docker Model Runner for local LLM execution.

## Two Methods for Model Runner Integration

This branch demonstrates two approaches to integrating with Docker Model Runner:

1. **Docker Compose** - Using the traditional Docker Compose approach (main branch)
2. **Testcontainers** - Using a programmatic approach with Testcontainers (this branch)

## Testcontainers Integration

This branch implements a Testcontainers-based approach for working with Docker Model Runner. It provides:

1. **Development Mode**: Run the application with automatic Model Runner connection setup
```bash
# Enable Model Runner in Docker Desktop with host-side TCP support first
go run -tags dev .
```

2. **Testing Mode**: Test connectivity to the Model Runner programmatically
```bash
cd tests
go test -v ./integration -run TestModelRunnerIntegration
```

### Benefits of the Testcontainers Approach

- **Dynamic Port Allocation**: Avoids port conflicts by using randomly assigned ports
- **Programmatic Control**: Start and manage connections via code rather than configuration files
- **Graceful Cleanup**: Automatic container cleanup when the application exits
- **Simplified Configuration**: Eliminates the need for multiple Docker Compose files
- **Unified Approach**: Same pattern for both development and testing

See `TESTCONTAINERS_USAGE.md` for detailed information about the Testcontainers implementation.

## Connecting to Model Runner

There are two ways to connect to Docker Model Runner:

### 1. Using Internal DNS

This method uses the internal Docker DNS resolution (`model-runner.docker.internal`)

### 2. Using TCP 

This method uses the host-side TCP support via `host.docker.internal:12434`

## Architecture

The application consists of three main components:

1. **Frontend**: React TypeScript application providing a responsive chat interface
2. **Backend**: Go server that handles API requests and connects to the LLM
3. **Model Runner**: Docker Model Runner running the Llama 3.2 (1B parameter) model

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │ >>> │   Backend   │ >>> │    Model    │
│  (React/TS) │     │    (Go)     │     │   Runner    │
└─────────────┘     └─────────────┘     └─────────────┘
      :3000              :8080              :12434
```

## Features

- Real-time chat interface with message history
- Streaming AI responses (tokens appear as they're generated)
- Dockerized deployment for easy setup
- Local LLM integration (no cloud API dependencies)
- Cross-origin resource sharing (CORS) enabled
- Comprehensive integration tests using Testcontainers

## Prerequisites

- Docker Desktop with Model Runner enabled
- Enable host-side TCP support in Model Runner settings
- Git
- Go 1.21 or higher

## Getting Started with Testcontainers

1. Clone this repository and check out this branch:
   ```bash
   git clone https://github.com/ajeetraina/genai-app-demo.git
   cd genai-app-demo
   git checkout feature/testcontainers-integration
   ```

2. Make sure Docker Model Runner is enabled in Docker Desktop

3. Run the application in development mode:
   ```bash
   go run -tags dev .
   ```

4. Access the frontend at [http://localhost:3000](http://localhost:3000)

## Testing Model Runner Integration

The application includes a dedicated test for Docker Model Runner connectivity:

```bash
cd tests
go test -v ./integration -run TestModelRunnerIntegration
```

This test will:

1. Start a socat container to forward traffic to Model Runner
2. Test connectivity to the Model Runner API
3. Display available models
4. Optionally attempt to create a new model
5. Provide the dynamically assigned URL for the application

## Additional Tests

```bash
# Run all tests
cd tests
go test -v ./integration

# Run specific test categories
go test -v ./integration -run TestGenAIAppIntegration    # API tests
go test -v ./integration -run TestFrontendIntegration    # UI tests
go test -v ./integration -run TestGenAIQuality           # Quality tests
go test -v ./integration -run TestGenAIPerformance       # Performance tests

# Run tests in short mode (faster)
go test -v ./integration -short
```

## Configuration

In the Testcontainers approach, environment variables are automatically set by the `dev_dependencies.go` file when running in development mode.

For manual configuration, you can set these environment variables:

- `BASE_URL`: URL for the model runner (set automatically in dev mode)
- `MODEL`: Model identifier to use (defaults to "ignaciolopezluna020/llama3.2:1B")
- `API_KEY`: API key for authentication (defaults to "ollama")

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
