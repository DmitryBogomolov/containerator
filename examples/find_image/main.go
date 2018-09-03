package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/DmitryBogomolov/containerator"
	"github.com/docker/docker/client"
)

func main() {
	repo := flag.String("repo", "", "Repo")
	tag := flag.String("tag", "", "Tag")
	flag.Parse()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	image, err := containerator.FindImage(cli, containerator.FindImageOptions{
		Repo: *repo,
		Tag:  *tag,
	})
	if err != nil {
		fmt.Println(err)
	} else if image == nil {
		fmt.Printf("image %s:%s is not found\n", *repo, *tag)
	} else {
		fmt.Printf("%s %s %v\n", image.ID[7:15], image.Name, time.Unix(image.Created, 0))
	}
}
