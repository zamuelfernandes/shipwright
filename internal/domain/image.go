package domain

import "context"

// Image representa uma imagem Docker salva localmente no host.
type Image struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    int64    `json:"size"`
	Created int64    `json:"created"`
}

// ImageRepository define o contrato para listagem de imagens Docker.
type ImageRepository interface {
	ListImages(ctx context.Context) ([]Image, error)
}
