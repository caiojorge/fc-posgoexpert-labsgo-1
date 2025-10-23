.PHONY: help build run test coverage clean docker-build docker-up docker-down deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/weather-api cmd/main.go

run: ## Run the application locally
	go run cmd/main.go

test: ## Run tests
	go test ./... -v

coverage: ## Generate test coverage report
	go test ./... -coverprofile=coverage.txt
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.txt coverage.html

docker-build: ## Build Docker image
	docker build -t weather-api .

docker-up: ## Start services with docker-compose
	docker compose up --build

docker-down: ## Stop services with docker-compose
	docker compose down

deploy: ## Deploy to Google Cloud Run (requires gcloud auth)
	gcloud builds submit --config cloudbuild.yaml

tidy: ## Tidy go modules
	go mod tidy

fmt: ## Format code
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

swagger: ## Generate swagger documentation
	swag init -g cmd/main.go
