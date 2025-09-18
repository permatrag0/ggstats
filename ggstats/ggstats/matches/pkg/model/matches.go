package model

type Tournament string
type Player1 string
type Player2 string
type Scorep1 int
type Scorep2 int
type Score int

type RecordID string

type RecordType string

type Matches struct {
	RecordID   string `json:"recordId"`
	RecordType string `json:"recordType"`
	Tournament string `json:"tournament"`
	Player1    string `json:"player1"`
	Player2    string `json:"player2"`
	Scorep1    int    `json:"scorep1"`
	Scorep2    int    `json:"scorep2"`
}
