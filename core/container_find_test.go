package core

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestFindContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContainers := []types.Container{
		{
			ID:      "00112233445566778899",
			Names:   []string{"/tester-1", "/tester-1a"},
			ImageID: "i1",
		},
		{
			ID:      "11223344556677889900",
			Names:   []string{},
			ImageID: "i2",
		},
		{
			ID:      "22334455667788990011",
			Names:   []string{"/tester-3"},
			ImageID: "i2",
		},
		{
			ID:      "33445566778899001122",
			Names:   []string{"/tester-4", "/tester-4a", "/tester-4b"},
			ImageID: "i1",
		},
		{
			ID:      "44556677889900112233",
			Names:   []string{},
			ImageID: "i2",
		},
	}

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(testContainers, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		cont, err := FindContainerByID(cli, "22334455667788990011")
		assert.NoError(t, err)
		assert.Equal(t, makeContainer(&testContainers[2]), cont)
	})

	t.Run("ByID / not found", func(t *testing.T) {
		cont, err := FindContainerByID(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, cont)
	})

	t.Run("ByShortID", func(t *testing.T) {
		cont, err := FindContainerByShortID(cli, "3344")
		assert.NoError(t, err)
		assert.Equal(t, makeContainer(&testContainers[3]), cont)
	})

	t.Run("ByShortID / not found", func(t *testing.T) {
		cont, err := FindContainerByID(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, cont)
	})

	t.Run("ByName", func(t *testing.T) {
		cont, err := FindContainerByName(cli, "tester-1")
		assert.NoError(t, err)
		assert.Equal(t, makeContainer(&testContainers[0]), cont)
	})

	t.Run("ByName / not found", func(t *testing.T) {
		cont, err := FindContainerByName(cli, "unknown")
		assert.NoError(t, err)
		assert.Nil(t, cont)
	})

	t.Run("ByImageID", func(t *testing.T) {
		conts, err := FindContainersByImageID(cli, "i2")
		assert.NoError(t, err)
		expected := []Container{
			makeContainer(&testContainers[1]),
			makeContainer(&testContainers[2]),
			makeContainer(&testContainers[4]),
		}
		assert.Equal(t, expected, conts)
	})

	t.Run("ByImageID / not found", func(t *testing.T) {
		conts, err := FindContainersByImageID(cli, "unknown")
		assert.NoError(t, err)
		assert.Equal(t, []Container(nil), conts)
	})
}
