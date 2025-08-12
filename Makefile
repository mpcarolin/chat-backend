.PHONY: help run test build clean dbuild up mock-run

# Default target
help:
	@echo "Available commands:"
	@echo "  make run        - Run the server with Azure provider"
	@echo "  make mock-run   - Run the server with mock provider"
	@echo "  make test       - Run all tests"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Remove built binaries"
	@echo "  make dbuild - Build Docker image"
	@echo "  make up   - Run with Docker Compose"

# Run the server with real chat provider
run:
	go run .

# Run the server with mock provider
mock-run:
	MOCK_CHAT=true go run .

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
