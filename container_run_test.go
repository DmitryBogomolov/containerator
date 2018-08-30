package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
)

func TestRunContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(container.ContainerCreateCreatedBody{ID: "id-1"}, nil)
	cli.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	cli.EXPECT().ContainerInspect(gomock.Any(), gomock.Any()).Return(types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    "id-1",
			Name:  "name-1",
			Image: "image-1",
			State: &types.ContainerState{
				Status: "status-1",
			},
		},
	}, nil)

	cont, err := RunContainer(cli, &ContainerOptions{
		Image: "image:1",
		Name:  "container-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := ContainerInfo{
		ID:    "id-1",
		Name:  "name-1",
		Image: "image-1",
		State: "status-1",
	}
	if *cont != expected {
		t.Fatalf("container: %v / expected: %v", cont, expected)
	}
}
