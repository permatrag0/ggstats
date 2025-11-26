package players

import (
	"context"
	"errors"

	metadataModel "ggstats.com/metadata/pkg"
	playersModel "ggstats.com/players/pkg/model"
	metadatapb "ggstats.com/proto/metadata"
)

var ErrNotFound = errors.New("players metadata not found")

// metadataGateway is what our controller expects from any metadata client.
type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatapb.Metadata, error)
}

// matchesGateway is kept here for future use; currently unused.
type matchesGateway interface {
	// e.g. later: GetMatches(ctx context.Context, recordID, recordType string) (...)
}

// Controller coordinates calls to metadata (and later matches)
// to build PlayerResults.
type Controller struct {
	matchesGateway  matchesGateway
	metadataGateway metadataGateway
}

// New constructs a new Controller.
// Even if matchesGateway is nil for now, we keep it in the signature
// so we can later plug matches in without breaking callers.
func New(matchesGateway matchesGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		matchesGateway:  matchesGateway,
		metadataGateway: metadataGateway,
	}
}

// Get builds a PlayerResults for the given id by calling the metadata service via gRPC.
func (c *Controller) Get(ctx context.Context, id string) (*playersModel.PlayerResults, error) {
	m, err := c.metadataGateway.Get(ctx, id)
	if err != nil {
		// TODO: inspect gRPC status and map NOT_FOUND to ErrNotFound if you want
		return nil, err
	}

	// Map gRPC metadata → the metadata model used inside PlayerResults
	md := metadataModel.Metadata{
		ID:       m.GetId(),
		Gamertag: m.GetGamertag(),
		Region:   m.GetRegion(),
		Sponsor:  m.GetSponsor(),
	}

	return &playersModel.PlayerResults{
		Metadata: md,
		// Scorep1 / Scorep2 left nil for now until matches is wired
	}, nil
}
