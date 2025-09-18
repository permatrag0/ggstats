package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"ggstats.com/matches/internal/controller/matches"
	httpHandler "ggstats.com/matches/internal/handler/http"
	"ggstats.com/matches/internal/repository/memory"
	"ggstats.com/pkg/discovery/consul"
	discovery "ggstats.com/pkg/registry"
)

const serviceName = "matches"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting matches service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
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
	repo := memory.New()
	ctrl := matches.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/matches", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
