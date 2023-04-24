// Program find_image shows usage of *containerator* functions that find docker images.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DmitryBogomolov/containerator/core"
	"github.com/docker/docker/client"
)

func displayImage(image core.Image) string {
	return fmt.Sprintf("%s (%s)", image.FullName(), image.ShortID())
}

func findImageByID(cli *client.Client, id string) error {
	image, err := core.FindImageByShortID(cli, id)
	if _, ok := err.(*core.ImageNotFoundError); ok {
		fmt.Println("Image not found")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println(displayImage(image))
	return nil
}

func findImageByName(cli *client.Client, name string) error {
	image, err := core.FindImageByName(cli, name)
	if _, ok := err.(*core.ImageNotFoundError); ok {
		fmt.Println("Image not found")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println(displayImage(image))
	return nil
}

func findAllImagesByName(cli *client.Client, name string) error {
	images, err := core.FindAllImagesByName(cli, name)
	if err != nil {
		return err
	}
	for _, image := range images {
		fmt.Println(displayImage(image))
	}
	return nil
}

func listAllImages(cli *client.Client) error {
	imageIDs, err := core.ListAllImageIDs(cli)
	if err != nil {
		return err
	}
	for _, id := range imageIDs {
		fmt.Println(id)
	}
	return nil
}

func run() error {
	var id string
	flag.StringVar(&id, "id", "", "image by id")
	var name string
	flag.StringVar(&name, "name", "", "image by name")
	var baseName string
	flag.StringVar(&baseName, "base-name", "", "all images by name")
	var listAll bool
	flag.BoolVar(&listAll, "list-all", false, "list all images")

	flag.Parse()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if id != "" {
		return findImageByID(cli, id)
	} else if name != "" {
		return findImageByName(cli, name)
	} else if baseName != "" {
		return findAllImagesByName(cli, baseName)
	} else if listAll {
		return listAllImages(cli)
	} else {
		flag.Usage()
		return nil
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
