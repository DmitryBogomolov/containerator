package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestRemoveContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerRemove(gomock.Any(), "0123456789ab", types.ContainerRemoveOptions{Force: true}).Return(nil)

	err := RemoveContainer(cli, testContainer("0123456789ab", ""))
	assert.NoError(t, err)
}
