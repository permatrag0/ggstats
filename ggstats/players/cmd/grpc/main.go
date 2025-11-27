package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	grpcmatches "ggstats.com/players/internal/gateway/matches/grpc"
	grpcmetadata "ggstats.com/players/internal/gateway/metadata/grpc"
	playerspb "ggstats.com/proto/players"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type playersServer struct {
	playerspb.UnimplementedPlayersServiceServer
	metadataClient *grpcmetadata.Client
	matchesClient  *grpcmatches.Client
}

func newPlayersServer(metadataClient *grpcmetadata.Client, matchesClient *grpcmatches.Client) *playersServer {
	return &playersServer{
		metadataClient: metadataClient,
		matchesClient:  matchesClient,
	}
}

func (s *playersServer) GetPlayer(ctx context.Context, req *playerspb.GetPlayerRequest) (*playerspb.GetPlayerResponse, error) {
	id := req.GetId()

	// 1) Get metadata
	md, err := s.metadataClient.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2) Get matches (we assume record_type="player")
	matches, err := s.matchesClient.GetMatches(ctx, id)
	if err != nil {
		// you could log and continue instead of failing hard if you prefer
		return nil, err
	}

	// 3) Map matches to players service proto
	var pmatches []*playerspb.Match
	for _, m := range matches {
		pmatches = append(pmatches, &playerspb.Match{
			RecordId:   m.GetRecordId(),
			RecordType: m.GetRecordType(),
			Tournament: m.GetTournament(),
			Player1:    m.GetPlayer1(),
			Player2:    m.GetPlayer2(),
			Scorep1:    m.GetScorep1(),
			Scorep2:    m.GetScorep2(),
		})
	}

	player := &playerspb.Player{
		Id:       md.GetId(),
		Gamertag: md.GetGamertag(),
		Region:   md.GetRegion(),
		Sponsor:  md.GetSponsor(),
		Matches:  pmatches,
	}

	return &playerspb.GetPlayerResponse{
		Player: player,
	}, nil
}

func main() {
	port := flag.String("port", "50053", "gRPC listen port")
	flag.Parse()

	// Where are the other services?
	metadataAddr := getEnv("METADATA_GRPC_ADDR", "localhost:50051")
	matchesAddr := getEnv("MATCHES_GRPC_ADDR", "localhost:50052")

	metadataClient, err := grpcmetadata.New(metadataAddr)
	if err != nil {
		log.Fatalf("failed to create metadata client: %v", err)
	}
	defer metadataClient.Close()

	matchesClient, err := grpcmatches.New(matchesAddr)
	if err != nil {
		log.Fatalf("failed to create matches client: %v", err)
	}
	defer matchesClient.Close()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	playerspb.RegisterPlayersServiceServer(grpcServer, newPlayersServer(metadataClient, matchesClient))

	reflection.Register(grpcServer)

	log.Printf("Players gRPC server listening on :%s", *port)
	log.Printf("Using metadata at %s, matches at %s", metadataAddr, matchesAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
