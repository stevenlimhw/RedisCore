
run: build
	@./bin/rediscore

build:
	@go build -o bin/rediscore