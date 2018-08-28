package containerator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type imageDesc struct {
	tag     string
	created int64
}

type imageDescList []*imageDesc

func (list imageDescList) Len() int {
	return len(list)
}

func (list imageDescList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list imageDescList) Less(i, j int) bool {
	return list[i].created > list[j].created
}

func selectImageDesc(baseTag string, image *types.ImageSummary) *imageDesc {
	for _, tag := range image.RepoTags {
		if strings.HasPrefix(tag, baseTag) {
			return &imageDesc{tag, image.Created}
		}
	}
	return nil
}

func filterImagesByRepoTag(baseTag string, images []types.ImageSummary) string {
	var descList imageDescList
	for _, image := range images {
		desc := selectImageDesc(baseTag, &image)
		if desc != nil {
			descList = append(descList, desc)
		}
	}
	if len(descList) == 0 {
		return ""
	}
	sort.Sort(descList)
	return descList[0].tag
}

func findImageRepoTag(cli *client.Client, baseTag string) (string, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return "", err
	}

	tag := filterImagesByRepoTag(baseTag, images)
	if tag == "" {
		return "", fmt.Errorf("image %s is not found", baseTag)
	}
	return tag, nil
}
