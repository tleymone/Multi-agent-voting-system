package ballotagt

import (
	"sync"
	"td3/comsoc"
	"time"
)

type BallotAgent struct {
	sync.Mutex
	ID         string
	ReqCount   int
	Addr       string
	Profile    comsoc.Profile
	NbAlts     int
	VoterIds   []string
	VotingFunc string
	Deadline   time.Time
}

func NewBallotAgent(id, addr string, nbAlts int, v []string, f string, d time.Time) *BallotAgent {
	return &BallotAgent{ID: id, Addr: addr, VoterIds: v, VotingFunc: f, Deadline: d}
}

func (ba *BallotAgent) CheckIsExistingVoter(agentId string) bool {
	for _, agt := range ba.VoterIds {
		if agt == agentId {
			return true
		}
	}
	return false
}

func (ba *BallotAgent) CheckPref(pref []comsoc.Alternative) bool {
	if ba.NbAlts != len(pref) {
		return false
	}
	for _, alt := range pref {
		if alt < 0 || alt > comsoc.Alternative(ba.NbAlts) {
			return false
		}
	}
	return true
}

func (ba *BallotAgent) DeadlineEnded() bool {
	return time.Now().After(ba.Deadline)
}
