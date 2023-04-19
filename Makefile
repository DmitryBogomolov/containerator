install-lint:
	go install golang.org/x/lint/golint@latest

install-mockgen:
	go install github.com/golang/mock/mockgen@latest

generate:
	go generate ./...

test:
	go test -v -coverprofile=coverage.out ./...

lint:
	golint ./...

build:
	go build -v ./...

view-cover:
	go tool cover -html=coverage.out
