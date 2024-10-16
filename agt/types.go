package agt

import (
	"td3/comsoc"
)

type NewBallotRequest struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	VoterIds []string `json:"voter-ids"`
	NbAlts   int      `json:"#alts"`
}

type NewBallotResponse struct {
	BallotId string `json:"ballot-id"`
}

type VoteRequest struct {
	AgentId string               `json:"agent-id"`
	VoteId  string               `json:"vote-id"`
	Prefs   []comsoc.Alternative `json:"prefs"`
	Options []int                `json:"options"`
}

type ResultRequest struct {
	BallotId string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  int   `json:"winner"`
	Ranking []int `json:"ranking"`
}
