package containerator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ImageInfo contains image information.
type ImageInfo struct {
	ID      string
	Tag     string
	Created int64
}

type imageInfoList []*ImageInfo

func (list imageInfoList) Len() int {
	return len(list)
}

func (list imageInfoList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list imageInfoList) Less(i, j int) bool {
	return list[i].Created > list[j].Created
}

func selectImageInfo(tagPrefix string, image *types.ImageSummary) *ImageInfo {
	for _, tag := range image.RepoTags {
		if strings.HasPrefix(tag, tagPrefix) {
			return &ImageInfo{
				ID:      image.ID,
				Tag:     tag,
				Created: image.Created,
			}
		}
	}
	return nil
}

func filterImagesByTag(tagPrefix string, images []types.ImageSummary) *ImageInfo {
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

// FindImageByTag selects images by tag prefix.
func FindImageByTag(cli client.ImageAPIClient, tagPrefix string) (*ImageInfo, error) {
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
