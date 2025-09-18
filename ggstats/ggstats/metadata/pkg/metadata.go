package model

type Metadata struct {
	ID       string `json:"id"`
	Gamertag string `json:"gamertag"`
	Region   string `json:"region"`
	Sponsor  string `json:"sponsor"`
}
