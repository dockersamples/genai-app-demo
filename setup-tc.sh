cat setup-tc.sh
#!/bin/bash
# setup-testcontainers.sh


# Run the necessary Go commands
go mod download && go mod tidy


# Optionally, run tests to verify everything works
# Uncomment if you want this to happen automatically
# go test ./...

cd tests/
go mod download && go mod tidy
go test -v ./integration -run TestModelRunnerIntegration
