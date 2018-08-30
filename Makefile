install-mockgen:
	go get -v github.com/golang/mock/mockgen

generate-mocks:
	mockgen -destination test_mocks/imageapiclient.go -package test_mocks github.com/docker/docker/client ImageAPIClient
	mockgen -destination test_mocks/containerapiclient.go -package test_mocks github.com/docker/docker/client ContainerAPIClient

test:
	go test -v
