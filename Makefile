.PHONY: build test lint clean run docker-up docker-down docker-logs

# Build the mata binary
build:
	go build -o bin/mata ./cmd/mata

# Run tests with coverage
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests for a specific package
test-pkg:
	go test -v -race ./$(PKG)/...

# Lint the codebase
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Run the application
run:
	go run ./cmd/mata $(ARGS)

# Install development dependencies
install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Initialize go module dependencies
init:
	go mod tidy
	go mod download

# Docker Compose Commands
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f