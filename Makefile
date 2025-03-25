# Variables
BINARY_NAME=myapp
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Commands
default: build

build: ## Compile the application
	go build -o $(BINARY_NAME) cmd/main.go

run: build ## Build and run the application
	./$(BINARY_NAME)

test: ## Run tests
	go test ./...

lint: ## Run linter
	golangci-lint run

fmt: ## Format the code
	go fmt ./...

clean: ## Remove built files
	rm -f $(BINARY_NAME)

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
