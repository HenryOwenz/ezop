# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=cloudgate

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test coverage deps build-linux build-windows build-all lint

all: lint test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 -v

build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_arm64 -v

build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 -v

build-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_darwin_arm64 -v

build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64.exe -v

build-windows-arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_windows_arm64.exe -v

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 build-windows-arm64

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
	rm -f release/*
	rm -rf release

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

lint:
	golangci-lint run

deps:
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy

# Release builds all binaries and creates a release
release: clean build-all
	@echo "Creating release..."
	@mkdir -p release
	@mv $(BINARY_NAME)_linux_* release/
	@mv $(BINARY_NAME)_darwin_* release/
	@mv $(BINARY_NAME)_windows_* release/
	@echo "Release binaries created in release/"
	@echo "Built for the following platforms:"
	@ls -l release/

# Run builds and executes the binary
run: build
	./$(BINARY_NAME)

# Install builds and installs the binary
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/cg 

# Test targets
.PHONY: test test-verbose test-coverage test-integration test-unit

# Main test target that runs all tests
test:
	$(GOTEST) -v ./...

test-verbose:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

test-integration:
	$(GOTEST) -v ./internal/ui/integration_test

test-unit:
	$(GOTEST) -v ./internal/ui/model ./internal/ui/update ./internal/ui/view
