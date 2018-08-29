package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("tag is not defined"))
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	image, err := containerator.FindImageByTag(cli, os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s %s %v\n", image.ID[7:15], image.Tag, time.Unix(image.Created, 0))
	}
}
