
.PHONY = all test

APP_NAME=rinha

all: test build

test:
	@echo "Running tests..."
	go test -v ./... -cover

build:
	@echo "Building '$(APP_NAME)'..."
	go clean
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(APP_NAME) cmd/rinha/main.go

swagger:
	@echo "Generating swagger..."
	swag init -g internal/http/api.go
	swag fmt

docker-build:
	@echo "Building docker image..."
	docker build -t marcusadriano/rinha-de-backend-2024-q1-go:latest .