install-mockgen:
	go get -v github.com/golang/mock/mockgen

generate-mocks:
	mockgen -destination test_mocks/imageapiclient.go -package test_mocks github.com/docker/docker/client ImageAPIClient
	mockgen -destination test_mocks/containerapiclient.go -package test_mocks github.com/docker/docker/client ContainerAPIClient

install:
	go get -t ./...

test:
	go test -v -coverprofile=coverage.out ./...

view-cover:
	go tool cover -html=coverage.out
