# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTIDY=$(GOCMD) mod tidy
GOCLEAN=$(GOCMD) clean
BINARY_NAME=ddg-search
MAIN_PATH=.

# Linting
GOLINT=golangci-lint

.PHONY: all build test clean tidy test-ci lint dev

all: tidy lint build test


build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

test-ci:
	$(GOTEST) -v ./... -coverprofile=coverage.out
	$(GOCMD) tool cover -func=coverage.out

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

tidy:
	$(GOTIDY)

lint:
	$(GOLINT) run

lint-fix:
	$(GOLINT) run --fix

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Docker commands
docker-build:
	docker build -t ddg-search:latest .

docker-run:
	docker run -p 8080:8080 ddg-search:latest

dev:
	@echo "Running in development mode with gowatch..."
	@gowatch -n -e ".go" -command "go run $(MAIN_PATH)/main.go"
