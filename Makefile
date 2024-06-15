dev:
	@go run cmd/main.go

run: build
	@./bin/bot

build:
	@go build -o bin/bot cmd/main.go