package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/DmitryBogomolov/containerator/core"
	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
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
