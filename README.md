# GenAI Application Testing 

This branch specifically checks the integration between your application and the model runner. It performs several key testing operations:

- Environment Setup Testing: It verifies that Testcontainers could successfully set up the required environment, including creating and configuring a socat container for networking.
- API Endpoint Testing: It tested various endpoints of the model runner service:

   - Checked the /models endpoint (to list available models)
   - Tested the /engines endpoint (to verify the service was running)

- Model Creation Testing: It attemptes to create a model (as an optional test) and verified this operation was successful.
- Connectivity Testing: The test verified that the model runner was accessible at a specific URL and provided the correct configuration URL to use in the application.

This type of testing is crucial for GenAI applications because it ensures that:

- The infrastructure components (containers, networking) work correctly
- The model service is properly configured and responding
- Basic model operations (like model loading/creation) function as expected

The approach you've implemented with Testcontainers is particularly valuable for GenAI testing because it provides an isolated, reproducible environment that can be used across development, testing, and CI/CD workflows, ensuring consistent behavior of your GenAI application regardless of where it runs.


## Getting Started

```
chmod +x setup-tc.sh
```


## Running the script

```
./setup-tc.sh
```

