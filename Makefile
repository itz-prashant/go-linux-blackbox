.PHONY: build run

build:
	@go build -o bin/main ./cmd/blackbox

run: build
	@./bin/main