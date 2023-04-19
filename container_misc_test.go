package containerator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
)

func TestRemoveContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerRemove(gomock.Any(), "container-1", types.ContainerRemoveOptions{Force: true}).Return(nil)

	err := RemoveContainer(cli, "container-1")
	assert.NoError(t, err)
}

func TestSuspendContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerRename(gomock.Any(), "my-id", gomock.Any()).Return(nil)
	cli.EXPECT().ContainerStop(gomock.Any(), "my-id", container.StopOptions{}).Return(nil)

	err := SuspendContainer(cli, "my-id")
	assert.NoError(t, err)
}

func TestResumeContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerRename(gomock.Any(), "my-id", "my-container").Return(nil)
	cli.EXPECT().ContainerStart(gomock.Any(), "my-id", types.ContainerStartOptions{}).Return(nil)

	err := ResumeContainer(cli, "my-id", "my-container")
	assert.NoError(t, err)
}
