package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestGetContainerName(t *testing.T) {
	t.Run("Empty name", func(t *testing.T) {
		name := GetContainerName(&types.Container{Names: []string{}})
		assert.Equal(t, "", name)
	})

	t.Run("First name", func(t *testing.T) {
		name := GetContainerName(&types.Container{Names: []string{"/c1", "/c2"}})
		assert.Equal(t, "c1", name)
	})
}

func TestGetContainerShortID(t *testing.T) {
	id := GetContainerShortID(&types.Container{ID: "01234567890123456789"})
	assert.Equal(t, "012345678901", id)
}

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
		assert.Equal(t, &testContainers[2], cont)
	})

	t.Run("ByID / not found", func(t *testing.T) {
		cont, err := FindContainerByID(cli, "unknown")
		assert.Error(t, err)
		contErr, ok := err.(*ContainerNotFoundError)
		assert.True(t, ok && contErr.Container == "unknown")
		assert.Nil(t, cont)
	})

	t.Run("ByShortID", func(t *testing.T) {
		cont, err := FindContainerByShortID(cli, "3344")
		assert.NoError(t, err)
		assert.Equal(t, &testContainers[3], cont)
	})

	t.Run("ByShortID / not found", func(t *testing.T) {
		cont, err := FindContainerByID(cli, "unknown")
		assert.Error(t, err)
		contErr, ok := err.(*ContainerNotFoundError)
		assert.True(t, ok && contErr.Container == "unknown")
		assert.Nil(t, cont)
	})

	t.Run("ByName", func(t *testing.T) {
		cont, err := FindContainerByName(cli, "tester-1")
		assert.NoError(t, err)
		assert.Equal(t, &testContainers[0], cont)
	})

	t.Run("ByName / not found", func(t *testing.T) {
		cont, err := FindContainerByName(cli, "unknown")
		assert.Error(t, err)
		contErr, ok := err.(*ContainerNotFoundError)
		assert.True(t, ok && contErr.Container == "unknown")
		assert.Nil(t, cont)
	})

	t.Run("ByImageID", func(t *testing.T) {
		conts, err := FindContainersByImageID(cli, "i2")
		assert.NoError(t, err)
		expected := []*types.Container{&testContainers[1], &testContainers[2], &testContainers[4]}
		assert.Equal(t, expected, conts)
	})

	t.Run("ByImageID / not found", func(t *testing.T) {
		conts, err := FindContainersByImageID(cli, "unknown")
		assert.NoError(t, err)
		var expected []*types.Container
		assert.Equal(t, expected, conts)
	})
}
