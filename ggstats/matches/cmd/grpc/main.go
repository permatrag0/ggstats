package main

import (
	"context"
	"flag"
	"log"
	"net"

	"ggstats.com/matches/internal/controller/matches"
	memoryrepo "ggstats.com/matches/internal/repository/memory"
	matchesmodel "ggstats.com/matches/pkg/model"
	matchespb "ggstats.com/proto/matches"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type matchesServer struct {
	matchespb.UnimplementedMatchesServiceServer
	ctrl *matches.Controller
}

func newMatchesServer(ctrl *matches.Controller) *matchesServer {
	return &matchesServer{ctrl: ctrl}
}

func (s *matchesServer) GetMatches(ctx context.Context, req *matchespb.GetMatchesRequest) (*matchespb.GetMatchesResponse, error) {
	recordID := matchesmodel.RecordID(req.GetRecordId())
	recordType := matchesmodel.RecordType(req.GetRecordType())

	list, err := s.ctrl.GetMatches(ctx, recordID, recordType)
	if err != nil {
		return nil, err
	}

	resp := &matchespb.GetMatchesResponse{}
	for _, m := range list {
		resp.Matches = append(resp.Matches, &matchespb.Match{
			RecordId:   m.RecordID,
			RecordType: m.RecordType,
			Tournament: m.Tournament,
			Player1:    m.Player1,
			Player2:    m.Player2,
			Scorep1:    int32(m.Scorep1),
			Scorep2:    int32(m.Scorep2),
		})
	}
	return resp, nil
}

func (s *matchesServer) PutMatch(ctx context.Context, req *matchespb.PutMatchRequest) (*matchespb.PutMatchResponse, error) {
	m := req.GetMatch()
	recordID := matchesmodel.RecordID(m.GetRecordId())
	recordType := matchesmodel.RecordType(m.GetRecordType())

	match := &matchesmodel.Matches{
		RecordID:   m.GetRecordId(),
		RecordType: m.GetRecordType(),
		Tournament: m.GetTournament(),
		Player1:    m.GetPlayer1(),
		Player2:    m.GetPlayer2(),
		Scorep1:    int(m.GetScorep1()),
		Scorep2:    int(m.GetScorep2()),
	}

	if err := s.ctrl.PutMatches(ctx, recordID, recordType, match); err != nil {
		return nil, err
	}
	return &matchespb.PutMatchResponse{}, nil
}

func main() {
	port := flag.String("port", "50052", "gRPC listen port")
	flag.Parse()

	repo := memoryrepo.New()
	ctrl := matches.New(repo)

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchespb.RegisterMatchesServiceServer(grpcServer, newMatchesServer(ctrl))

	reflection.Register(grpcServer)

	log.Printf("Matches gRPC server listening on :%s", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
