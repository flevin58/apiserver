run:
	@go run .

build:
	@go build -o bin/apiserver .

test:
	@go test -v ./...
