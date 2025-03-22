# GenAI Application Testing 

This GitHub branch specifically checks the integration between your application and the model runner. It performs several key testing operations:

- Environment Setup Testing: It verifies that Testcontainers could successfully set up the required environment, including creating and configuring a socat container for networking.
- API Endpoint Testing: It tests various endpoints of the model runner service:
   - Checks the /models endpoint (to list available models)
   - Tests the /engines endpoint (to verify the service was running)
- Model Creation Testing: It attempts to create a model (as an optional test) and verified this operation was successful.
- Connectivity Testing: The test verified that the model runner was accessible at a specific URL and provided the correct configuration URL to use in the application.

This type of testing is crucial for GenAI applications because it ensures that:

- The infrastructure components (containers, networking) work correctly
- The model service is properly configured and responding
- Basic model operations (like model loading/creation) function as expected

Why testcontainers? The approach implemented with Testcontainers is particularly valuable for GenAI testing because it provides an isolated, reproducible environment that can be used across development, testing, and CI/CD workflows, ensuring consistent behavior of your GenAI application regardless of where it runs.


## Getting Started

```
chmod +x setup-tc.sh
```


## Running the script

```
./setup-tc.sh
```

You will see the following result:

```
./setup-tc.sh
=== RUN   TestModelRunnerIntegration
2025/03/22 11:30:01 github.com/testcontainers/testcontainers-go - Connected to docker:
  Server Version: 28.0.2 (via Testcontainers Desktop 1.18.1)
  API Version: 1.43
  Operating System: Docker Desktop
  Total Memory: 9937 MB
  Resolved Docker Host: tcp://127.0.0.1:49295
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 2cfa4d78c19d284c11acd58ccc7793aefda24b4dccf93701d9895a3d6e080814
  Test ProcessID: 6b7d8317-5c2a-42e3-a1a2-5a490caafa13
2025/03/22 11:30:01 🐳 Creating container for image testcontainers/ryuk:0.6.0
2025/03/22 11:30:01 ✅ Container created: c4b4712ee09f
2025/03/22 11:30:01 🐳 Starting container: c4b4712ee09f
2025/03/22 11:30:01 ✅ Container started: c4b4712ee09f
2025/03/22 11:30:01 🚧 Waiting for container id c4b4712ee09f image: testcontainers/ryuk:0.6.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2025/03/22 11:30:01 🐳 Creating container for image alpine/socat
2025/03/22 11:30:01 ✅ Container created: 89fae4dea4e4
2025/03/22 11:30:01 🐳 Starting container: 89fae4dea4e4
2025/03/22 11:30:01 ✅ Container started: 89fae4dea4e4
2025/03/22 11:30:01 🚧 Waiting for container id 89fae4dea4e4 image: alpine/socat. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
    model_runner_test.go:64: Testing GET /models endpoint...
    model_runner_test.go:85: Available models: 0 found
    model_runner_test.go:96: Testing /engines endpoint...
    model_runner_test.go:108: Engines endpoint response: Docker Model Runner

        The service is running.
    model_runner_test.go:130: Attempting to create model (optional test)...
    model_runner_test.go:148: Successfully created model
    model_runner_test.go:159: Model Runner test completed successfully!
    model_runner_test.go:160: Model Runner is accessible via: http://127.0.0.1:55536
    model_runner_test.go:161: Use this URL in your application config: http://127.0.0.1:55536/engines/llama.cpp/v1
2025/03/22 11:30:08 🐳 Terminating container: 89fae4dea4e4
2025/03/22 11:30:09 🚫 Container terminated: 89fae4dea4e4
--- PASS: TestModelRunnerIntegration (7.80s)
PASS
ok  	github.com/ajeetraina/genai-app-demo/tests/integration	(cached)
```



Here's a breakdown of what's happening:

- Test Initialization:
  - The test connects to your Docker environment (Docker Desktop version 28.0.2)
  - It reports memory (9937 MB) and connection details
  - A unique test session ID is generated to track this test run

- Container Setup:
   - Ryuk Container: First, a container using the testcontainers/ryuk:0.6.0 image is created and started. Ryuk is a cleanup service that ensures containers are removed when tests complete, even if they terminate unexpectedly.
   - Socat Container: Then, an alpine/socat container is created and started. This container acts as a network proxy, forwarding traffic between your tests and the model runner service.

- Testing Process:
    - Models Endpoint: The test checks the /models endpoint and reports that 0 models were found.
    - Engines Endpoint: The test verifies the /engines endpoint is responding with "Docker Model Runner The service is running."
    - Model Creation: The test attempts to create a model and reports success.


- Test Completion:
   - The test passes successfully (in 7.80 seconds)
   - The socat container is terminated
   - The test provides important URL information:
   - Model Runner is accessible via: http://127.0.0.1:55536
   - Use this URL in your application config: http://127.0.0.1:55536/engines/llama.cpp/v1

