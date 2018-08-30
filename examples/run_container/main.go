package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func main() {
	imagePtr := flag.String("image", "", "image name")
	namePtr := flag.String("name", "", "container name")
	volumePtr := flag.String("volume", "", "volume")
	portPtr := flag.String("port", "", "port")
	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := &containerator.ContainerOptions{
		Image: *imagePtr,
		Name:  *namePtr,
	}
	if *volumePtr != "" {
		list := strings.Split(*volumePtr, ":")
		options.Volumes = make(map[string]string)
		options.Volumes[list[0]] = list[1]
	}
	if *portPtr != "" {
		list := strings.Split(*portPtr, ":")
		options.Ports = make(map[int]int)
		from, _ := strconv.Atoi(list[0])
		to, _ := strconv.Atoi(list[1])
		options.Ports[from] = to
	}

	cont, err := containerator.RunContainer(cli, options)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", cont.Image[7:15], cont.ID[7:15], cont.Name, cont.State)
}
