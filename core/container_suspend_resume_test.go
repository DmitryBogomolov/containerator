package core_test

import (
	"testing"

	. "github.com/DmitryBogomolov/containerator/core"
	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
