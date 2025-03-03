.DEFAULT_GOAL := test

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-more
lint-more:
	golangci-lint run --enable gocognit,cyclop,funlen,gocyclo ./...


.PHONY: test
test:
	go test -cover ./...

.PHONY: race
race:
	go test -cover -race ./...
