package grpcmatches

import (
	"context"

	matchespb "ggstats.com/proto/matches"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client matchespb.MatchesServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: matchespb.NewMatchesServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// GetMatches for a specific player (we treat playerID as record_id and "player" as record_type)
func (c *Client) GetMatches(ctx context.Context, playerID string) ([]*matchespb.Match, error) {
	resp, err := c.client.GetMatches(ctx, &matchespb.GetMatchesRequest{
		RecordId:   playerID,
		RecordType: "player",
	})
	if err != nil {
		return nil, err
	}
	return resp.GetMatches(), nil
}
