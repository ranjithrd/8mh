.PHONY: install dev build run clean test

install:
	go mod download
	go mod tidy

dev:
	@command -v air > /dev/null || go install github.com/air-verse/air@latest
	air

build:
	go build -o bin/server src/main.go

run:
	go run src/main.go

clean:
	rm -rf bin/ tmp/

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

swagger:
	swag init -g src/main.go -o src/docs

swagger-fmt:
	swag fmt -g src/main.go
