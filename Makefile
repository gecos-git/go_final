build:
	@go build -o bin/todo cmd/main.go

test:
	@go test ./tests

run:
	@./bin/todo
