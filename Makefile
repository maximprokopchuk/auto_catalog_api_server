.PHONY: build

build: 
		go build -v ./cmd/auto_catalog_api_server

.PHONY: run

run: 
		go run -v ./cmd/auto_catalog_api_server

.PHONY: test
test:
		go test -v -race -timeout 30s ./...

.PHONY: deps
deps:
		go mod tidy

.PHONY: lint
lint:
		golangci-lint run


DEFAULT_GOAL := build