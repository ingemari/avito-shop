up:
	@docker compose up -d

down:
	@docker compose down -v

go:
	@go run cmd/main.go

test:
	@go test ./internal/services -cover

.PHONY: up down go test
