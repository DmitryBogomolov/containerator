package containerator

import (
	"errors"
	"os"
	"strings"
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
		expectedContainer := types.Container{ID: "cid1"}
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{expectedContainer}, nil)

		cont, err := RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
		})

		assertEqual(t, err, nil, "error")
		assertEqual(t, cont.ID, expectedContainer.ID, "container")
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
		expectedErr := errors.New("error-on-start")
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(expectedErr)
		cli.EXPECT().
			ContainerRemove(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)

		_, err := RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
		})

		assertEqual(t, err, expectedErr, "error")
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
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				types.Container{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Volumes: []Mapping{
				Mapping{"/src1", "/dst1"},
				Mapping{"/src2", "/dst2"},
			},
			Ports: []Mapping{
				Mapping{"1001", "1000"},
				Mapping{"2001", "2000"},
			},
		})
	})

	t.Run("Env", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		os.Setenv("D", "test")
		defer os.Unsetenv("D")
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{
					Image: "image:1",
					Env: []string{
						"A=1",
						"B=",
						"C=3",
						"D=test",
					},
				},
				&container.HostConfig{},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				types.Container{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Env: []Mapping{
				Mapping{"A", "1"},
				Mapping{"B", ""},
				Mapping{"C", "3"},
				Mapping{"D", ""},
			},
		})
	})

	t.Run("EnvReader", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{
					Image: "image:1",
					Env: []string{
						"A=0",
						"A=1",
						"B=2",
					},
				},
				&container.HostConfig{},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				types.Container{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Env: []Mapping{
				Mapping{"A", "1"},
				Mapping{"B", "2"},
			},
			EnvReader: strings.NewReader("#test\nA=0\n"),
		})
	})

	t.Run("RestartPolicy", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{Image: "image:1"},
				&container.HostConfig{
					RestartPolicy: container.RestartPolicy{
						Name: "on-failure",
					},
				},
				nil, "container-1").
			Return(container.ContainerCreateCreatedBody{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				types.Container{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image:         "image:1",
			Name:          "container-1",
			RestartPolicy: RestartOnFailure,
		})
	})
}
