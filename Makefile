.PHONY: build test test-all test-integration test-coverage dev-up dev-down run-ingester run-exporter

build:
	go build ./...

test:
	go test ./... -count=1

test-verbose:
	go test ./... -count=1 -v

test-all: test test-integration

test-integration: dev-up
	@echo "Waiting for MinIO to be ready..."
	@until curl -sf http://localhost:9000/minio/health/live > /dev/null 2>&1; do sleep 1; done
	go test -tags=integration ./internal/integration/ -v -count=1

test-coverage:
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo ""
	@echo "HTML report: go tool cover -html=coverage.out"

dev-up:
	docker compose up -d

dev-down:
	docker compose down -v

run-ingester:
	go run ./cmd/ingester config.local.json

run-exporter:
	go run ./cmd/exporter config.local.json
