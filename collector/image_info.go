package collector

import (
	"strings"
)

type ImageInfo struct {
	BaseImage string
	Tag       string
	Hash      string
}

func parseDockerImage(image string) ImageInfo {
	var info ImageInfo

	parts := strings.SplitN(image, "@", 2)

	if len(parts) > 1 {
		info.Hash = parts[1]
	}

	imageParts := strings.SplitN(parts[0], ":", 2)

	info.BaseImage = imageParts[0]

	if len(imageParts) > 1 {
		info.Tag = imageParts[1]
	}

	return info
}
