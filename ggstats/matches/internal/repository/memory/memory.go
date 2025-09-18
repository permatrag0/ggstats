package memory

import (
	"context"

	"ggstats.com/matches/internal/repository"
	"ggstats.com/matches/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Matches
}

func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Matches{}}
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Matches, error) {

	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if matches, ok := r.data[recordType][recordID]; !ok || len(matches) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, matches *model.Matches) error {

	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Matches{}
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *matches)
	return nil
}
