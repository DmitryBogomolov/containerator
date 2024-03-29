package core

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"

	"github.com/DmitryBogomolov/containerator/test_mocks"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/golang/mock/gomock"
	"gopkg.in/yaml.v2"
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
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
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

		assert.NoError(t, err)
		assert.Equal(t, expectedContainer.ID, cont.ID())
	})

	t.Run("RemoveNonStarted", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{Image: "image:1"},
				&container.HostConfig{},
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
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

		assert.Equal(t, expectedErr, err)
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
						"1000/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "1001"}},
						"2000/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "2001"}},
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: "/src1",
							Target: "/dst1",
						},
						{
							Type:   mount.TypeBind,
							Source: "/src2",
							Target: "/dst2",
						},
					},
				},
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Volumes: []Mapping{
				{"/src1", "/dst1"},
				{"/src2", "/dst2"},
			},
			Ports: []Mapping{
				{"1001", "1000"},
				{"2001", "2000"},
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
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Env: []Mapping{
				{"A", "1"},
				{"B", ""},
				{"C", "3"},
				{"D", ""},
			},
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
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image:         "image:1",
			Name:          "container-1",
			RestartPolicy: container.RestartPolicyOnFailure,
		})
	})

	t.Run("Network", func(t *testing.T) {
		cli := test_mocks.NewMockContainerAPIClient(ctrl)
		cli.EXPECT().
			ContainerCreate(
				gomock.Any(),
				&container.Config{Image: "image:1"},
				&container.HostConfig{
					NetworkMode: container.NetworkMode("test-net"),
				},
				nil, nil, "container-1").
			Return(container.CreateResponse{ID: "cid1"}, nil)
		cli.EXPECT().
			ContainerStart(gomock.Any(), "cid1", gomock.Any()).
			Return(nil)
		cli.EXPECT().
			ContainerList(gomock.Any(), gomock.Any()).
			Return([]types.Container{
				{},
			}, nil)

		RunContainer(cli, &RunContainerOptions{
			Image:   "image:1",
			Name:    "container-1",
			Network: "test-net",
		})
	})
}

func TestRunContainerOptions(t *testing.T) {
	t.Run("JSONMarshal", func(t *testing.T) {
		options := RunContainerOptions{
			Image:   "image:1",
			Name:    "container-1",
			Network: "network-1",
			Volumes: []Mapping{
				{"/src1", "/dst1"},
			},
			Env: []Mapping{
				{"A", "1"},
				{"B", "2"},
			},
		}

		bytes, err := json.MarshalIndent(options, "", "  ")
		data := string(bytes)

		assert.NoError(t, err)
		expected := strings.Join([]string{
			`{`,
			`  "image": "image:1",`,
			`  "name": "container-1",`,
			`  "volumes": [`,
			`    {`,
			`      "/src1": "/dst1"`,
			`    }`,
			`  ],`,
			`  "env": [`,
			`    {`,
			`      "A": "1"`,
			`    },`,
			`    {`,
			`      "B": "2"`,
			`    }`,
			`  ],`,
			`  "network": "network-1"`,
			`}`,
		}, "\n")
		assert.Equal(t, expected, data)
	})

	t.Run("JSONUnmarshal", func(t *testing.T) {
		data := strings.Join([]string{
			`{`,
			`"image": "image:1",`,
			`"name": "container-1",`,
			`"env": [`,
			`{"A": "1" },`,
			`{ "B": "2" }`,
			`]`,
			`}`,
		}, "\n")
		var options RunContainerOptions

		err := json.Unmarshal([]byte(data), &options)

		assert.NoError(t, err)
		assert.Equal(t, RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Env:   []Mapping{{"A", "1"}, {"B", "2"}},
		}, options)
	})

	t.Run("YAMLMarshal", func(t *testing.T) {
		options := RunContainerOptions{
			Image:   "image:1",
			Name:    "container-1",
			Network: "network-1",
			Volumes: []Mapping{
				{"/src1", "/dst1"},
			},
			Env: []Mapping{
				{"A", "1"},
				{"B", "2"},
			},
		}

		bytes, err := yaml.Marshal(options)
		data := string(bytes)

		assert.NoError(t, err)
		expected := strings.Join([]string{
			"image: image:1",
			"name: container-1",
			"volumes:",
			"- /src1: /dst1",
			"env:",
			`- A: "1"`,
			`- B: "2"`,
			"network: network-1",
			"",
		}, "\n")
		assert.Equal(t, expected, data)
	})

	t.Run("YAMLUnmarshal", func(t *testing.T) {
		data := strings.Join([]string{
			"image: image:1",
			"name: container-1",
			"env:",
			"- A: 1",
			"- B: 2",
		}, "\n")
		var options RunContainerOptions

		err := yaml.Unmarshal([]byte(data), &options)

		assert.NoError(t, err)
		assert.Equal(t, RunContainerOptions{
			Image: "image:1",
			Name:  "container-1",
			Env:   []Mapping{{"A", "1"}, {"B", "2"}},
		}, options)
	})
}
