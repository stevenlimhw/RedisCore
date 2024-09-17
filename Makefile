
run: build
	@./bin/rediscore

build: fmt
	@go build -o bin/rediscore

fmt:
	@fieldalignment -fix ./...
	@go fmt ./...

test: lint
	@go test -v ./...

lint:
	@golangci-lint run
