.PHONY: build build-all install clean test fmt lint docker-up docker-down docker-logs docker-restart docker-clean

# Build for current platform
build:
	go build -o bin/outline cmd/outline/main.go

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/outline-linux-amd64 cmd/outline/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/outline-linux-arm64 cmd/outline/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/outline-darwin-amd64 cmd/outline/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/outline-darwin-arm64 cmd/outline/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/outline-windows-amd64.exe cmd/outline/main.go

# Install to local bin
install:
	go install cmd/outline/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Docker commands
docker-up:
	cd .docker-outline && docker compose up -d

docker-down:
	cd .docker-outline && docker compose down

docker-logs:
	cd .docker-outline && docker compose logs -f

docker-restart:
	cd .docker-outline && docker compose restart outline

docker-clean:
	cd .docker-outline && docker compose down -v
