.PHONY: test fmt lint run

test:
	go test -v -race -coverprofile=coverage.out ./...

fmt:
	golangci-lint run --fix

lint:
	golangci-lint run

run:
	go run cmd/server/main.go

