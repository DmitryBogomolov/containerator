package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestGetContainerName(t *testing.T) {
	var name string

	name = GetContainerName(&types.Container{Names: []string{}})
	assertEqual(t, name, "", "name")

	name = GetContainerName(&types.Container{Names: []string{"/c1", "/c2"}})
	assertEqual(t, name, "c1", "name")
}

func TestFindContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContainers := []types.Container{
		types.Container{
			ID:      "c1",
			Names:   []string{"/tester-1", "/tester-1a"},
			ImageID: "i1",
		},
		types.Container{
			ID:      "c2",
			Names:   []string{},
			ImageID: "i2",
		},
		types.Container{
			ID:      "c3",
			Names:   []string{"/tester-3"},
			ImageID: "i2",
		},
		types.Container{
			ID:      "c4",
			Names:   []string{"/tester-4", "/tester-4a", "/tester-4b"},
			ImageID: "i1",
		},
		types.Container{
			ID:      "c5",
			Names:   []string{},
			ImageID: "i2",
		},
	}

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(testContainers, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		var cont *types.Container
		var err error

		cont, err = FindContainerByID(cli, "c3")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, &testContainers[2], "container")

		cont, err = FindContainerByID(cli, "c5")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, &testContainers[4], "container")

		cont, err = FindContainerByID(cli, "unknown")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, nil, "container")
	})

	t.Run("ByName", func(t *testing.T) {
		var cont *types.Container
		var err error

		cont, err = FindContainerByName(cli, "tester-1")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, &testContainers[0], "container")

		cont, err = FindContainerByName(cli, "tester-4a")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, &testContainers[3], "container")

		cont, err = FindContainerByName(cli, "unknown")
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, nil, "container")
	})

	t.Run("ByImageID", func(t *testing.T) {
		var conts []*types.Container
		var err error

		conts, err = FindContainersByImageID(cli, "i2")
		assertEqual(t, err, nil, "error")
		assertEqual(t, len(conts), 3, "containers count")
		assertEqual(t, conts[0], &testContainers[1], "container 1")
		assertEqual(t, conts[1], &testContainers[2], "container 2")
		assertEqual(t, conts[2], &testContainers[4], "container 3")

		conts, err = FindContainersByImageID(cli, "unknown")
		assertEqual(t, err, nil, "error")
		assertEqual(t, len(conts), 0, "containers count")
	})
}
