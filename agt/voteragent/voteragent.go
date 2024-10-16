package voteragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"td3/agt"
	"td3/comsoc"
)

var nbAgent int
var nbVotes = 0

type VoterAgent struct {
	ID      string
	Url     string
	Name    string
	Prefs   []comsoc.Alternative
	Options []int
}

type AgentI interface {
	Equal(ag VoterAgent) bool
	DeepEqual(ag VoterAgent) bool
	Clone() VoterAgent
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) (bool, error)
	Start()
}

func NewVoterAgent(name string, url string, prefs []comsoc.Alternative, o []int) *VoterAgent {
	nbAgent++
	return &VoterAgent{ID: fmt.Sprintf("ag_id%01d", nbAgent), Url: "http://localhost:" + url, Name: name, Prefs: prefs, Options: o}
}

func (va *VoterAgent) Equal(ag VoterAgent) bool {
	return va.ID == ag.ID
}

func (va *VoterAgent) DeepEqual(ag VoterAgent) bool {
	for i, v := range va.Prefs {
		if v != ag.Prefs[i] {
			return false
		}
	}
	return (va.ID == ag.ID) && (va.Name == ag.Name) && (va.Url == ag.Url)
}

func (va *VoterAgent) Clone() *VoterAgent {
	return NewVoterAgent(va.Name, va.Url, va.Prefs, va.Options)
}

func (va *VoterAgent) String() string {
	return va.Name
}

func (va *VoterAgent) Prefers(alt1 comsoc.Alternative, alt2 comsoc.Alternative) (bool, error) {
	//vérifier que les alternatives sont bien dans la préférence
	if comsoc.Rank(alt1, va.Prefs) == -1 || comsoc.Rank(alt2, va.Prefs) == -1 {
		return false, errors.New("one of the alternatives is not in the preference list of this agent")
	}
	return comsoc.IsPref(alt1, alt2, va.Prefs), nil
}

func (va *VoterAgent) doRequest(ballot string) (err error) {
	req := agt.VoteRequest{
		AgentId: va.ID,
		VoteId:  ballot,
		Prefs:   va.Prefs,
		Options: va.Options,
	}
	nbVotes++
	// sérialisation de la requête
	url := va.Url + "/vote"
	data, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	return
}

func (va *VoterAgent) Start(ballot string) {
	err := va.doRequest(ballot)

	if err != nil {
		log.Fatal(va.ID, "error:", err.Error())
	} else {
		log.Printf("%s a envoyé %d à %s. (Options : %d)", va.ID, va.Prefs, ballot, va.Options)
	}
}
