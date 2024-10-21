package collector

import (
	"reflect"
	"testing"
)

func TestParseDockerImage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ImageInfo
	}{
		{
			name:  "Full image with base, tag, and hash",
			input: "ghcr.io/synapto-platform/n8n:build_3@sha256:5cb3b7ef347a590dec108dc3cb0532853dc57751a8d42ed7dcabdd113e31373f",
			expected: ImageInfo{
				BaseImage: "ghcr.io/synapto-platform/n8n",
				Tag:       "build_3",
				Hash:      "sha256:5cb3b7ef347a590dec108dc3cb0532853dc57751a8d42ed7dcabdd113e31373f",
			},
		},
		{
			name:  "Image with base and tag, no hash",
			input: "ubuntu:20.04",
			expected: ImageInfo{
				BaseImage: "ubuntu",
				Tag:       "20.04",
				Hash:      "",
			},
		},
		{
			name:  "Image with base and hash, no tag",
			input: "nginx@sha256:4aabc3a08446fc3f7d74473e2bf7b6badef7129da3a0a187eb2f169f0e417a1b",
			expected: ImageInfo{
				BaseImage: "nginx",
				Tag:       "",
				Hash:      "sha256:4aabc3a08446fc3f7d74473e2bf7b6badef7129da3a0a187eb2f169f0e417a1b",
			},
		},
		{
			name:  "Image with only base",
			input: "alpine",
			expected: ImageInfo{
				BaseImage: "alpine",
				Tag:       "",
				Hash:      "",
			},
		},
		{
			name:  "Image with registry and base",
			input: "docker.io/library/redis",
			expected: ImageInfo{
				BaseImage: "docker.io/library/redis",
				Tag:       "",
				Hash:      "",
			},
		},
		{
			name:  "Image with registry, base, and latest tag",
			input: "docker.io/library/redis:latest",
			expected: ImageInfo{
				BaseImage: "docker.io/library/redis",
				Tag:       "latest",
				Hash:      "",
			},
		},
		{
			name:  "Empty string",
			input: "",
			expected: ImageInfo{
				BaseImage: "",
				Tag:       "",
				Hash:      "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDockerImage(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseDockerImage(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
