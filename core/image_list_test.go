package core

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAllImageIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testImages := []types.ImageSummary{
		{
			ID: "sha256:00112233445566778899",
		},
		{
			ID: "sha256:11223344556677889900",
		},
		{
			ID: "sha256:22334455667788990011",
		},
		{
			ID: "sha256:33445566778899001122",
		},
	}

	cli := test_mocks.NewMockImageAPIClient(ctrl)
	cli.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(testImages, nil).AnyTimes()

	imageIDs, err := ListAllImageIDs(cli)

	assert.NoError(t, err)
	assert.Equal(t, []string{
		"sha256:00112233445566778899",
		"sha256:11223344556677889900",
		"sha256:22334455667788990011",
		"sha256:33445566778899001122",
	}, imageIDs)
}
