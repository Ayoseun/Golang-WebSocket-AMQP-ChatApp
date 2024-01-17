BINARY_NAME = golang-websocket-amqp-chatapp
.DEFAULT_GOAL := run
.PHONY: build run clean dep vet lint

# Create target builds
build:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/$(BINARY_NAME) -buildvcs=false ./cmd/main.go
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/$(BINARY_NAME)-darwin -buildvcs=false ./cmd/main.go

# Run command
run: build
#	@./bin/$(BINARY_NAME)-darwin
	@./bin/${BINARY_NAME}

# Clean package
clean:
	@go clean
	@rm ./bin/$(BINARY_NAME)
	@rm ./bin/$(BINARY_NAME)-darwin

# Download packages
dep:
	@go mod download
#reporting suspicious constructs 
vet:
	@go vet

lint:
	@golangci-lint run --enable-all
