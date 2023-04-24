package core

import (
	"testing"

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
	cli.EXPECT().ContainerRename(gomock.Any(), "0123456789ab", gomock.Any()).Return(nil)
	cli.EXPECT().ContainerStop(gomock.Any(), "0123456789ab", container.StopOptions{}).Return(nil)

	err := SuspendContainer(cli, testContainer("0123456789ab", ""))
	assert.NoError(t, err)
}

func TestResumeContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerRename(gomock.Any(), "0123456789ab", "my-container").Return(nil)
	cli.EXPECT().ContainerStart(gomock.Any(), "0123456789ab", types.ContainerStartOptions{}).Return(nil)

	err := ResumeContainer(cli, testContainer("0123456789ab", ""), "my-container")
	assert.NoError(t, err)
}
