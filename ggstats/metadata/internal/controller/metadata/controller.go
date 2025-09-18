package metadata

import (
	"context"
	"errors"

	model "ggstats.com/metadata/pkg"
)

var ErrNotFound = errors.New("Not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, metadata *model.Metadata) error
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)

	if err != nil {
		return nil, ErrNotFound
	}

	return res, err
}

func (c *Controller) Create(ctx context.Context, metadata *model.Metadata) error {
	return c.repo.Put(ctx, metadata)
}
