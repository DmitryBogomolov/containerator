package core_test

//go:generate mockgen -destination ./test_mocks/mock_imageapiclient.go -package test_mocks github.com/docker/docker/client ImageAPIClient
//go:generate mockgen -destination ./test_mocks/mock_containerapiclient.go -package test_mocks github.com/docker/docker/client ContainerAPIClient
