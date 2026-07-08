package domain

import "context"

// Image represents a local Docker image.
type Image struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    int64    `json:"size"`
	Created int64    `json:"created"`
}

// ImageRepository defines the contract for listing Docker images.
type ImageRepository interface {
	ListImages(ctx context.Context) ([]Image, error)
}
