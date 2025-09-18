package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ggstats.com/metadata/internal/controller/metadata"
	httphandler "ggstats.com/metadata/internal/handler/http"
	"ggstats.com/metadata/internal/repository/memory"
	"ggstats.com/pkg/discovery/consul"
	discovery "ggstats.com/pkg/registry"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting metadata service on port %d", port)
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
	r := memory.New()
	c := metadata.New(r)
	h := httphandler.New(c)

	http.Handle("/metadata", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetMetadata(w, r)
		case http.MethodPost:
			h.CreateMetadata(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
