package http

import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"ggstats.com/matches/pkg/model"
	discovery "ggstats.com/pkg/registry"
	// "ggstats.com/players/internal/gateway"
)

type Gateway struct {
	registry discovery.Registry
}

// putRating implements movie.ratingGateway.
func (g *Gateway) PutMatch(ctx context.Context, recordID model.RecordID, recordType model.RecordType, matches *model.Matches) error {
	addrs, err := g.registry.ServiceAddress(ctx, "matches")
	if err != nil {
		return err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/matches"
	log.Printf("%s", "Calling matches service, request: PUT"+url)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("tournament", fmt.Sprintf("%v", matches.Tournament))
	values.Add("PLayer1", fmt.Sprintf("%v", matches.Player1))
	values.Add("Player2", fmt.Sprintf("%v", matches.Player2))
	values.Add("Scorep1", fmt.Sprintf("%v", matches.Scorep1))
	values.Add("Scorep2", fmt.Sprintf("%v", matches.Scorep2))

	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non 2xx response :%v", resp)
	}
	return nil
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
// 	addrs, err := g.registry.ServiceAddress(ctx, "rating")
// 	if err != nil {
// 		return 0, err
// 	}
// 	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
// 	log.Printf("%s", "Calling rating service, request: GET "+url)
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return 0, err
// 	}
// 	req = req.WithContext(ctx)
// 	values := req.URL.Query()
// 	values.Add("id", string(recordID))
// 	values.Add("type", fmt.Sprintf("%v", recordType))
// 	req.URL.RawQuery = values.Encode()
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode == http.StatusNotFound {
// 		return 0, gateway.ErrNotFound
// 	} else if resp.StatusCode/100 != 2 {
// 		return 0, fmt.Errorf("non-2xx response: %v", resp)
// 	}
// 	var v float64
// 	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
// 		return 0, err
// 	}
// 	return v, nil
// }
