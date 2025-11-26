package main

import (
	"log"
	"net/http"
	"os"

	"ggstats.com/players/internal/controller/players"
	grpcmetadata "ggstats.com/players/internal/gateway/metadata/grpc"
	httphandler "ggstats.com/players/internal/handler/http"
)

func main() {
	// Where is the metadata gRPC server?
	// In local dev, we'll default to localhost:50051.
	metadataAddr := os.Getenv("METADATA_GRPC_ADDR")
	if metadataAddr == "" {
		metadataAddr = "localhost:50051"
	}

	metadataClient, err := grpcmetadata.New(metadataAddr)
	if err != nil {
		log.Fatalf("failed to create metadata gRPC client: %v", err)
	}
	defer func() {
		if err := metadataClient.Close(); err != nil {
			log.Printf("failed to close metadata client: %v", err)
		}
	}()

	// We don't use matches yet, so pass nil for the matchesGateway.
	ctrl := players.New(nil, metadataClient)
	h := httphandler.New(ctrl)

	http.Handle("/players", http.HandlerFunc(h.GetPlayersDetails))

	addr := ":8080"
	log.Printf("Players HTTP API listening on %s (metadata gRPC at %s)", addr, metadataAddr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
