package containerator

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func TestRunContainer(t *testing.T) {
	stub := func(args ...interface{}) (interface{}, error) {
		if len(args) == 5 {
			return container.ContainerCreateCreatedBody{ID: "id-1"}, nil
		}
		if len(args) == 3 {
			return nil, nil
		}
		if len(args) == 2 {
			return types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:    "id-1",
					Name:  "name-1",
					Image: "image-1",
					State: &types.ContainerState{
						Status: "status-1",
					},
				},
			}, nil
		}
		t.Fatalf("unexpected stub call\n")
		return nil, nil
	}

	cli := &testContainerAPIClient{stub: stub}
	cont, err := RunContainer(cli, &ContainerOptions{
		Image: "image:1",
		Name:  "container-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := ContainerInfo{
		ID:    "id-1",
		Name:  "name-1",
		Image: "image-1",
		State: "status-1",
	}
	if *cont != expected {
		t.Fatalf("container: %v / expected: %v", cont, expected)
	}
}
