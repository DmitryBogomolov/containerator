package core

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAllContainerIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContainers := []types.Container{
		{
			ID: "00112233445566778899",
		},
		{
			ID: "11223344556677889900",
		},
		{
			ID: "22334455667788990011",
		},
		{
			ID: "33445566778899001122",
		},
		{
			ID: "44556677889900112233",
		},
	}

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(testContainers, nil).AnyTimes()

	containerIDs, err := ListAllContainerIDs(cli)

	assert.NoError(t, err)
	assert.Equal(t, []string{
		"00112233445566778899",
		"11223344556677889900",
		"22334455667788990011",
		"33445566778899001122",
		"44556677889900112233",
	}, containerIDs)
}
