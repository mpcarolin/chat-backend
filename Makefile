.PHONY: help run test build clean dbuild up watch run-mock run-azure run-ollama

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "Provider-specific commands:"
	@echo "  make run-mock   - Run the server with mock provider (default)"
	@echo "  make run-azure  - Run the server with Azure Q&A provider"
	@echo "  make run-ollama - Run the server with Ollama provider"
	@echo ""
	@echo "General commands:"
	@echo "  make run        - Run the server with default provider"
	@echo "  make test       - Run all tests"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Remove built binaries"
	@echo ""
	@echo "Docker commands:"
	@echo "  make dbuild     - Build Docker image"
	@echo "  make up         - Run with Docker Compose"
	@echo "  make watch      - Run with Docker Compose watch (auto-rebuild)"
	@echo ""
	@echo "Environment variables for Azure:"
	@echo "  AZURE_QNA_ENDPOINT, AZURE_QNA_API_KEY, AZURE_QNA_PROJECT_NAME, AZURE_QNA_DEPLOYMENT_NAME"
	@echo ""
	@echo "Environment variables for Ollama:"
	@echo "  OLLAMA_BASE_URL (default: http://localhost:11434)"
	@echo "  OLLAMA_MODEL (default: mistral)"

# Run the server with default provider (mock)
run:
	go run .

# Run the server with mock provider
run-mock:
	CHAT_PROVIDER=mock go run .

# Run the server with Azure Q&A provider
run-azure:
	@if [ -z "$$AZURE_QNA_ENDPOINT" ] || [ -z "$$AZURE_QNA_API_KEY" ] || [ -z "$$AZURE_QNA_PROJECT_NAME" ] || [ -z "$$AZURE_QNA_DEPLOYMENT_NAME" ]; then \
		echo "Error: Azure environment variables not set."; \
		echo "Please set: AZURE_QNA_ENDPOINT, AZURE_QNA_API_KEY, AZURE_QNA_PROJECT_NAME, AZURE_QNA_DEPLOYMENT_NAME"; \
		exit 1; \
	fi
	CHAT_PROVIDER=azure-qa go run .

# Run the server with Ollama provider
run-ollama:
	@echo "Starting server with Ollama provider..."
	@echo "Using OLLAMA_BASE_URL: $${OLLAMA_BASE_URL:-http://localhost:11434}"
	@echo "Using OLLAMA_MODEL: $${OLLAMA_MODEL:-mistral}"
	CHAT_PROVIDER=ollama go run .

# Run all tests
test:
	go test -v ./...

# Build the application
build:
	go build -o chat-backend

# Clean built binaries
clean:
	rm -f chat-backend

# Build Docker image
dbuild:
	docker-compose build

# Run with Docker Compose
up:
	docker-compose up

# Run with Docker Compose watch (auto-rebuild on file changes)
watch:
	docker-compose watch
