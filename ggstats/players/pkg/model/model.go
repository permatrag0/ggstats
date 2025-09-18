package model

import model "ggstats.com/metadata/pkg"

type PlayerResults struct {
	Scorep1  *int           `json:"scorep1,omitempty"`
	Scorep2  *int           `json:"scorep2,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
