package containerator

import (
	"errors"
	"testing"

	"github.com/docker/go-connections/nat"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/golang/mock/gomock"
)

func TestRunContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("CreateAndRun", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{Image: "image:1"},
				&container.HostConfig{},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerInspect(gomock.Any(), "cid1").
			Return(types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:    "id-1",
					Name:  "name-1",
					Image: "image-1",
					State: &types.ContainerState{
						Status: "status-1",
					},
				},
			}, nil)

		cont, err := RunContainer(cli, &ContainerOptions{
			Image: "image:1",
			Name:  "container-1",
		})

		assertEqual(t, err, nil, "error")
		assertEqual(t, *cont, ContainerInfo{
			ID:    "id-1",
			Name:  "name-1",
			Image: "image-1",
			State: "status-1",
		}, "container")
	})

	t.Run("VolumesAndPorts", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		var dummy struct{}
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{
					Image: "image:1",
					ExposedPorts: nat.PortSet{
						"1000/tcp": dummy,
						"2000/tcp": dummy,
					},
				},
				&container.HostConfig{
					PortBindings: nat.PortMap{
						"1000/tcp": []nat.PortBinding{nat.PortBinding{HostIP: "0.0.0.0", HostPort: "1001"}},
						"2000/tcp": []nat.PortBinding{nat.PortBinding{HostIP: "0.0.0.0", HostPort: "2001"}},
					},
					Mounts: []mount.Mount{
						mount.Mount{
							Type:   mount.TypeBind,
							Source: "/src1",
							Target: "/dst1",
						},
						mount.Mount{
							Type:   mount.TypeBind,
							Source: "/src2",
							Target: "/dst2",
						},
					},
				},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerInspect(gomock.Any(), "cid1").
			Return(types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{},
				},
			}, nil)

		RunContainer(cli, &ContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Volumes: map[string]string{
				"/src1": "/dst1",
				"/src2": "/dst2",
			},
			Ports: map[int]int{
				1000: 1001,
				2000: 2001,
			},
		})
	})

	t.Run("RemoveNonStarted", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{Image: "image:1"},
				&container.HostConfig{},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		expected := errors.New("error-on-start")
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(expected)
		cli.EXPECT().
			ContainerRemove(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)

		_, err := RunContainer(cli, &ContainerOptions{
			Image: "image:1",
			Name:  "container-1",
		})

		assertEqual(t, err, expected, "error")
	})
}
