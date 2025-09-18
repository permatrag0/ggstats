package memory

import (
	"context"
	"sync"

	"ggstats.com/metadata/internal/repository"
	model "ggstats.com/metadata/pkg"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

func (r *Repository) Put(ctx context.Context, metadata *model.Metadata) error {
	r.data[metadata.ID] = metadata
	return nil
}
