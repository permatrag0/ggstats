package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ggstats.com/pkg/discovery/consul"
	discovery "ggstats.com/pkg/registry"
	"ggstats.com/players/internal/controller/players"
	matchesgateway "ggstats.com/players/internal/gateway/matches/http"
	metadatagateway "ggstats.com/players/internal/gateway/metadata/http"
	httphandler "ggstats.com/players/internal/handler/http"
)

const serviceName = "players"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API Handler port")
	flag.Parse()
	log.Printf("Starting players + match service on port: %d", port)
	registry, err := consul.NewRegistry(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadatagateway.New(registry)
	matchesGateway := matchesgateway.New(registry)
	ctrl := players.New(matchesGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/players", http.HandlerFunc(h.GetPlayersDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
