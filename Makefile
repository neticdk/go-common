.DEFAULT_GOAL := test

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: race
race:
	go test -cover -race ./...
