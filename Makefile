install-lint:
	go install golang.org/x/lint/golint@latest

install-mockgen:
	go install github.com/golang/mock/mockgen@latest

generate-mocks:
	mockgen -destination test_mocks/imageapiclient.go -package test_mocks github.com/docker/docker/client ImageAPIClient
	mockgen -destination test_mocks/containerapiclient.go -package test_mocks github.com/docker/docker/client ContainerAPIClient

test:
	go test -v -coverprofile=coverage.out ./...

lint:
	golint ./...

build:
	go build -v ./...

view-cover:
	go tool cover -html=coverage.out
