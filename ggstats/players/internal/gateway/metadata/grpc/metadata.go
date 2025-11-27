package grpcmetadata

import (
	"context"

	metadatapb "ggstats.com/proto/metadata"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client metadatapb.MetadataServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: metadatapb.NewMetadataServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Get(ctx context.Context, id string) (*metadatapb.Metadata, error) {
	resp, err := c.client.GetMetadata(ctx, &metadatapb.GetMetadataRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetMetadata(), nil
}
