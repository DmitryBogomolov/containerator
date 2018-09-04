package containerator

import (
	"testing"

	"github.com/DmitryBogomolov/containerator/test_mocks"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestFindContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContainers := []types.Container{
		types.Container{
			ID:      "c1",
			Names:   []string{"/tester-1", "/tester-1a"},
			ImageID: "i1",
			Image:   "image-1",
			State:   "running",
		},
		types.Container{
			ID:      "c2",
			Names:   []string{},
			ImageID: "i2",
			Image:   "image-2",
			State:   "stopped",
		},
		types.Container{
			ID:      "c3",
			Names:   []string{"/tester-3"},
			ImageID: "i2",
			Image:   "image-2a",
			State:   "exited",
		},
		types.Container{
			ID:      "c4",
			Names:   []string{"/tester-4", "/tester-4a", "/tester-4b"},
			ImageID: "i1",
			Image:   "image-1",
			State:   "running",
		},
		types.Container{
			ID:      "c5",
			Names:   []string{},
			ImageID: "i2",
			Image:   "image-2b",
			State:   "testing",
		},
	}

	cli := test_mocks.NewMockContainerAPIClient(ctrl)
	cli.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(testContainers, nil).AnyTimes()

	t.Run("ByID", func(t *testing.T) {
		var cont *ContainerInfo
		var err error

		cont, err = FindContainer(cli, FindContainerOptions{ID: "c3"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *cont, ContainerInfo{ID: "c3", Image: "image-2a", Name: "tester-3", State: "exited"}, "container")

		cont, err = FindContainer(cli, FindContainerOptions{ID: "c5"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *cont, ContainerInfo{ID: "c5", Image: "image-2b", Name: "", State: "testing"}, "container")

		cont, err = FindContainer(cli, FindContainerOptions{ID: "unknown"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, nil, "container")
	})

	t.Run("ByName", func(t *testing.T) {
		var cont *ContainerInfo
		var err error

		cont, err = FindContainer(cli, FindContainerOptions{Name: "tester-1"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *cont, ContainerInfo{ID: "c1", Image: "image-1", Name: "tester-1", State: "running"}, "container")

		cont, err = FindContainer(cli, FindContainerOptions{Name: "tester-4a"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, *cont, ContainerInfo{ID: "c4", Image: "image-1", Name: "tester-4a", State: "running"}, "container")

		cont, err = FindContainer(cli, FindContainerOptions{Name: "unknown"})
		assertEqual(t, err, nil, "error")
		assertEqual(t, cont, nil, "container")
	})
}
