.PHONY: lint
pwd ?= $(shell pwd)

lint: 
	@docker run --rm \
		-v $(pwd):/app \
		-w /app \
		golangci/golangci-lint:v1.47.1 golangci-lint run -v

build:
	@docker build --tag nps-alerts .

run:
	@docker run \
    --env-file .env \
    --publish 8080:8080 \
    nps-alerts

test:
	go test ./src/... \
		-covermode=atomic \
		-timeout=10s \
		-race 