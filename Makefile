
APP_NAME=rinha

build:
	go clean
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(APP_NAME) cmd/rinha/main.go

swagger:
	swag init -g internal/http/api.go
	swag fmt

docker-build:
	docker build -t marcusadriano/rinha-de-backend-2024-q1-go:latest .