package players

import (
	"context"
	"errors"

	matchesModel "ggstats.com/matches/pkg/model"
	metadataModel "ggstats.com/metadata/pkg"
	"ggstats.com/players/internal/gateway"
	"ggstats.com/players/pkg/model"
)

var ErrNotFound = errors.New("players metadata not found")

type matchesGateway interface {
	// GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
	PutMatch(ctx context.Context, recordID matchesModel.RecordID, recordType matchesModel.RecordType, matches *matchesModel.Matches) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	matchesGateway  matchesGateway
	metadataGateway metadataGateway
}

func New(matchesGateway matchesGateway, metametadataGateway metadataGateway) *Controller {
	return &Controller{matchesGateway, metametadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.PlayerResults, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.PlayerResults{Metadata: *metadata}
	return details, nil
}
