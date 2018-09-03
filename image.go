package containerator

import (
	"errors"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ImageInfo contains image information.
type ImageInfo struct {
	ID      string
	Name    string
	Created int64
}

type imageSummaryListByCreated []*types.ImageSummary

func (list imageSummaryListByCreated) Len() int {
	return len(list)
}

func (list imageSummaryListByCreated) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list imageSummaryListByCreated) Less(i, j int) bool {
	return list[i].Created > list[j].Created
}

func getImageName(image *types.ImageSummary) string {
	if len(image.RepoTags) > 0 {
		return image.RepoTags[0]
	}
	return ""
}

func findImageByID(id string, images []types.ImageSummary) *ImageInfo {
	for _, image := range images {
		if image.ID == id {
			return &ImageInfo{
				ID:      id,
				Name:    getImageName(&image),
				Created: image.Created,
			}
		}
	}
	return nil
}

func filterImagesByRepo(repo string, images []types.ImageSummary) []*types.ImageSummary {
	var list []*types.ImageSummary
	for i, image := range images {
		for _, repoTag := range image.RepoTags {
			parts := strings.Split(repoTag, ":")
			if parts[0] == repo {
				list = append(list, &images[i])
				break
			}
		}
	}
	return list
}

func findNewestImage(images []*types.ImageSummary) *ImageInfo {
	if len(images) == 0 {
		return nil
	}
	sort.Sort(imageSummaryListByCreated(images))
	image := images[0]
	return &ImageInfo{
		ID:      image.ID,
		Name:    getImageName(image),
		Created: image.Created,
	}
}

func findImageByTag(tag string, images []*types.ImageSummary) *ImageInfo {
	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			parts := strings.Split(repoTag, ":")
			if len(parts) > 1 && parts[1] == tag {
				return &ImageInfo{
					ID:      image.ID,
					Name:    repoTag,
					Created: image.Created,
				}
			}
		}
	}
	return nil
}

// ErrFindImage shows that options are not valid.
var ErrFindImage = errors.New("neither *ID* nor *Repo* are provided")

// FindImageOptions defines search config.
type FindImageOptions struct {
	ID   string
	Repo string
	Tag  string
}

// FindImage selects images by tag prefix.
func FindImage(cli client.ImageAPIClient, options FindImageOptions) (*ImageInfo, error) {
	images, err := cliImageList(cli)
	if err != nil {
		return nil, err
	}

	if options.ID != "" {
		return findImageByID(options.ID, images), nil
	}
	if options.Repo != "" {
		list := filterImagesByRepo(options.Repo, images)
		if options.Tag != "" {
			return findImageByTag(options.Tag, list), nil
		}
		return findNewestImage(list), nil
	}
	return nil, ErrFindImage
}
