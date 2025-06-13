.PHONY: build clean test run dev install help

# Build configuration
BINARY_DIR := bin
VELO_BINARY := $(BINARY_DIR)/velo
VELOCTL_BINARY := $(BINARY_DIR)/veloctl

# Go configuration
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Default target
all: build

## build: Build all binaries
build: $(VELO_BINARY) $(VELOCTL_BINARY)
	@echo "✅ Build complete!"
	@echo "Binaries:"
	@echo "  🏗️  velo     - Server (manager/worker)"
	@echo "  🛠️  veloctl  - CLI client"

$(VELO_BINARY):
	@echo "🔨 Building velo server..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $@ ./cmd/velo

$(VELOCTL_BINARY):
	@echo "🔨 Building veloctl CLI..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $@ ./cmd/cli

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

## deps: Download and tidy dependencies
deps:
	@echo "📦 Managing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## run: Run manager server in development mode
run: $(VELO_BINARY)
	@echo "🚀 Starting Velo manager server..."
	@echo "📡 API: http://localhost:37355"
	@echo "🌐 Web: http://localhost:8080"
	@echo "👤 Default login: admin/admin"
	$(VELO_BINARY) --manager --web-port 8080

## dev: Run server with auto-restart (requires 'entr' tool)
dev:
	@if command -v entr >/dev/null 2>&1; then \
		echo "🔄 Running in development mode with auto-restart..."; \
		find . -name "*.go" | entr -r make run; \
	else \
		echo "❌ 'entr' tool not found. Install with:"; \
		echo "   Ubuntu/Debian: sudo apt install entr"; \
		echo "   macOS: brew install entr"; \
		echo "   Fedora: sudo dnf install entr"; \
		echo "   Arch: sudo pacman -S entr"; \
		exit 1; \
	fi

## install: Install binaries to system PATH
install: build
	@echo "📥 Installing binaries to /usr/local/bin..."
	sudo cp $(VELO_BINARY) /usr/local/bin/
	sudo cp $(VELOCTL_BINARY) /usr/local/bin/
	@echo "✅ Installation complete!"

## docker-build: Build Docker images
docker-build:
	@echo "🐳 Building Docker images..."
	docker build -t velo:latest .

## format: Format Go code
format:
	@echo "✨ Formatting code..."
	$(GOCMD) fmt ./...

## lint: Run linters (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "🔍 Running linters..."; \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install with:"; \
		echo "   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.50.1"; \
	fi

## release: Build release binaries for multiple platforms
release: clean
	@echo "🚀 Building release binaries..."
	@mkdir -p $(BINARY_DIR)/release
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/release/velo-linux-amd64 ./cmd/velo
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/release/veloctl-linux-amd64 ./cmd/cli
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/release/velo-linux-arm64 ./cmd/velo
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/release/veloctl-linux-arm64 ./cmd/cli
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/release/velo-darwin-amd64 ./cmd/velo
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/release/veloctl-darwin-amd64 ./cmd/cli
	# macOS arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/release/velo-darwin-arm64 ./cmd/velo
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_DIR)/release/veloctl-darwin-arm64 ./cmd/cli
	@echo "✅ Release binaries built in $(BINARY_DIR)/release/"

## help: Show this help
help:
	@echo "🚀 Velo Development Commands"
	@echo ""
	@echo "Available commands:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'