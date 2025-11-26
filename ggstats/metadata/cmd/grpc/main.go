// metadata/cmd/grpc/main.go
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"ggstats.com/metadata/internal/controller/metadata"
	memoryrepo "ggstats.com/metadata/internal/repository/memory"
	metadatamodel "ggstats.com/metadata/pkg"
	metadatapb "ggstats.com/proto/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type metadataServer struct {
	metadatapb.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func newMetadataServer(ctrl *metadata.Controller) *metadataServer {
	return &metadataServer{ctrl: ctrl}
}

func (s *metadataServer) GetMetadata(ctx context.Context, req *metadatapb.GetMetadataRequest) (*metadatapb.GetMetadataResponse, error) {
	m, err := s.ctrl.Get(ctx, req.GetId())
	if err != nil {
		// you can map your domain errors to gRPC status codes here
		return nil, err
	}

	return &metadatapb.GetMetadataResponse{
		Metadata: &metadatapb.Metadata{
			Id:       m.ID,
			Gamertag: m.Gamertag,
			Region:   m.Region,
			Sponsor:  m.Sponsor,
		},
	}, nil
}

func (s *metadataServer) CreateMetadata(ctx context.Context, req *metadatapb.CreateMetadataRequest) (*metadatapb.CreateMetadataResponse, error) {
	in := req.GetMetadata()
	m := &metadatamodel.Metadata{
		ID:       in.GetId(),
		Gamertag: in.GetGamertag(),
		Region:   in.GetRegion(),
		Sponsor:  in.GetSponsor(),
	}

	if err := s.ctrl.Create(ctx, m); err != nil {
		return nil, err
	}

	return &metadatapb.CreateMetadataResponse{
		Metadata: &metadatapb.Metadata{
			Id:       m.ID,
			Gamertag: m.Gamertag,
			Region:   m.Region,
			Sponsor:  m.Sponsor,
		},
	}, nil
}

func main() {
	port := flag.String("port", "50051", "gRPC listen port")
	flag.Parse()

	// Use your existing in-memory repository + controller
	repo := memoryrepo.New()
	ctrl := metadata.New(repo)

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	metadatapb.RegisterMetadataServiceServer(grpcServer, newMetadataServer(ctrl))

	// optional, but handy for local debugging with tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("Metadata gRPC server listening on %s", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
