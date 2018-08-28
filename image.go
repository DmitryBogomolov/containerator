package containerator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type imageInfo struct {
	id      string
	tag     string
	created int64
}

type imageInfoList []*imageInfo

func (list imageInfoList) Len() int {
	return len(list)
}

func (list imageInfoList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list imageInfoList) Less(i, j int) bool {
	return list[i].created > list[j].created
}

func selectImageInfo(tagPrefix string, image *types.ImageSummary) *imageInfo {
	for _, tag := range image.RepoTags {
		if strings.HasPrefix(tag, tagPrefix) {
			return &imageInfo{
				id:      image.ID,
				tag:     tag,
				created: image.Created,
			}
		}
	}
	return nil
}

func filterImagesByTag(tagPrefix string, images []types.ImageSummary) *imageInfo {
	var descList imageInfoList
	for _, image := range images {
		desc := selectImageInfo(tagPrefix, &image)
		if desc != nil {
			descList = append(descList, desc)
		}
	}
	if len(descList) == 0 {
		return nil
	}
	sort.Sort(descList)
	return descList[0]
}

func findImageByTag(cli client.ImageAPIClient, tagPrefix string) (*imageInfo, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}

	tag := filterImagesByTag(tagPrefix, images)
	if tag == nil {
		return nil, fmt.Errorf("image %s is not found", tagPrefix)
	}
	return tag, nil
}
