.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -v ./...

.PHONY: build
build:
	go build ./...
