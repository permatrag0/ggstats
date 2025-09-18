package matches

import (
	"context"
	"errors"

	"ggstats.com/matches/pkg/model"
)

var ErrNotFound = errors.New("matches not found for record")

// type Matches struct {
// 	RecordID   string `json:"recordid"`
// 	Tournament string `json:"tournament"`
// 	Player1    string `json:"player1"`
// 	Player2    string `json:"player2"`
// 	Scorep1    int    `json:"scorep1"`
// 	Scorep2    int    `json:"scorep2"`
// }

type matchesRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Matches, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, matches *model.Matches) error
}

type Controller struct {
	repo matchesRepository
}

func New(repo matchesRepository) *Controller {
	return &Controller{repo}
}

// func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
// 	ratings, err := c.repo.Get(ctx, recordID, recordType)
// 	if err != nil && err == repository.ErrNotFound {
// 		return 0, err
// 	} else if err != nil {
// 		return 0, err
// 	}

// 	sum := float64(0)
// 	for _, r := range ratings {
// 		sum += float64(r.Value)
// 	}

// 	return sum / float64(len(ratings)), nil
// }

func (c *Controller) GetMatches(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Matches, error) {
	matches, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil && errors.Is(err, ErrNotFound) { // Check for repository.ErrNotFound
		return nil, ErrNotFound // Return your controller's ErrNotFound
	} else if err != nil {
		return nil, err // Return other unexpected errors directly
	}
	return matches, nil
}

func (c *Controller) PutMatches(ctx context.Context, recordID model.RecordID, recordType model.RecordType, matches *model.Matches) error {
	return c.repo.Put(ctx, recordID, recordType, matches)
}
