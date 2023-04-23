[![CI](https://github.com/DmitryBogomolov/containerator/actions/workflows/ci.yml/badge.svg)](https://github.com/DmitryBogomolov/containerator/actions/workflows/ci.yml)

# containerator

## core

Helper functions to work with docker containers.

Several most common and basic operations are supported.

```go
// import "github.com/docker/docker/client"
cli, _ = client.NewEnvClient()

// docker image ls -aq
core.ListAllImageIDs(cli)

// docker ps -aq
core.ListAllContainerIDs(cli)

// docker inspect 8b5a55df88ec
core.FindImageByShortID(cli, "0123456789ab")

// docker inspect 8b5a55df88ec
core.FindContainerByShortID(cli, "0123456789ab")

// docker run -it -d --name my-container-1 --network my-network-1 -v /tmp:/usr/app -p 50001:3000 -e A=1 my-image:1
core.RunContainer(cli, &core.RunContainerOptions{
    Image: "my-image:1",
    Name: "my-container-1",
    RestartPolicy: RestartAlways,
    Network: "my-network-1",
    Volumes: []Mapping{
        {"/tmp", "/usr/app"},
    },
    Ports: []Mapping{
        {"50001", "3000"},
    },
    Env: []Mapping{
        {"A", "1"},
    },
})
```

## manage

Contains function to run, suspend, resume, remove containers.

### Examples

- [find_image](./examples/find_image/README.md)
- [find_container](./examples/find_container/README.md)
- [run_container](./examples/run_container/README.md)
- [manage_container](./examples/manage_container/README.md)
- [manage_container_server](./examples/manage_container_server/README.md)
