package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"td3/agt"
	ballotagt "td3/agt/ballotagent"
	"td3/comsoc"
	"time"
)

var url string // Port du serveur
var alreadyVotedMap = make(map[string][]string)
var nbBallot int
var ballotMap = make(map[string]*ballotagt.BallotAgent)
var ruleMap = make(map[string]func(p comsoc.Profile) (a comsoc.Alternative, e error))
var SWFMap = make(map[string]func(p comsoc.Profile) (a []comsoc.Alternative, e error))
var threeshold []int // Liste des seuils pour le vote approval

func decodeBallotRequest(r *http.Request) (req agt.NewBallotRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	parse, _ := time.Parse(time.RFC3339, req.Deadline) // le format est "2012-11-01T22:08:41+01:00"
	if req.Rule == "" || req.Deadline == "" || req.VoterIds == nil || req.NbAlts <= 0 {
		return agt.NewBallotRequest{}, errors.New("missing or unvalid field in request body")
	}
	if parse.Before(time.Now()) {
		return agt.NewBallotRequest{}, errors.New("deadline must be after current time")
	}
	if req.Rule != "majority" && req.Rule != "borda" && req.Rule != "copeland" && req.Rule != "kemeny" && req.Rule != "stv" && req.Rule != "approval" {
		return agt.NewBallotRequest{}, errors.New("rule must be one of ['majority', 'borda', 'copeland', 'kemeny', 'stv', 'approval']")
	}
	return
}

func decodeVoteRequest(r *http.Request) (req agt.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	if req.AgentId == "" || req.VoteId == "" || req.Prefs == nil || req.Options == nil {
		return agt.VoteRequest{}, errors.New("missing field in request body")
	}
	return
}

func decodeResultRequest(r *http.Request) (req agt.ResultRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	if req.BallotId == "" {
		return agt.ResultRequest{}, errors.New("missing field in request body")
	}
	return
}

func checkAlreadyVoted(ballotId, agentId string) bool {
	for _, agt := range alreadyVotedMap[ballotId] {
		if agt == agentId {
			return true
		}
	}
	return false
}

func checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func doNewBallot(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	req, err := decodeBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// création du ballot
	nbBallot++
	log.Println("Création du vote n°", nbBallot)
	parse, _ := time.Parse(time.RFC3339, req.Deadline)
	log.Println("  ↪  ", "ballot"+fmt.Sprint(nbBallot), ":", req.Rule, parse, req.VoterIds, req.NbAlts)

	// Le ballot s'appelle balloti (ex : ballot1 pour le premier)
	ballot := ballotagt.NewBallotAgent("ballot"+fmt.Sprint(nbBallot), url, req.NbAlts, req.VoterIds, req.Rule, parse)
	ballotMap[ballot.ID] = ballot

	// envoi de l'id du ballot
	var res agt.NewBallotResponse
	res.BallotId = ballot.ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	if req.Rule == "approval" {
		threeshold = nil
	}
}

func doVote(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	req, err := decodeVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	log.Println(req.AgentId, "a voté pour", req.VoteId, "avec les préférences", req.Prefs)

	// vérification que le ballot existe
	if ba, ok := ballotMap[req.VoteId]; ok {
		// mise à jour du nombre de requêtes
		ba.Lock()
		defer ba.Unlock()
		ba.ReqCount++

		// vérification de la méthode de la requête
		if !checkMethod("POST", w, r) {
			return
		}

		//vérifier que le votant fait partie du ballot et qu'il a donc le droit de voter
		if !ba.CheckIsExistingVoter(req.AgentId) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "this agent does not exist in the ballot")
			return
		}

		//Vérifie si le ballot est un vote par approbation
		if ba.VotingFunc == "approval" {
			if len(req.Options) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "the threshold of acceptance is needed")
				return
			}
			threeshold = append(threeshold, req.Options[0])
		}

		//vérification que la préférence est valide
		if ba.CheckPref(req.Prefs) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "preference is invalid")
			return
		}

		//vérification que l'agent n'a pas déjà voté
		if checkAlreadyVoted(ba.ID, req.AgentId) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, errors.New("this agent has already voted and cannot vote twice"))
			return
		}

		//vérification de la deadline
		if ba.DeadlineEnded() {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, errors.New("ballot has already reached its deadline"))
			return
		}

		// traitement de la requête
		ba.Profile = append(ba.Profile, req.Prefs)
		alreadyVotedMap[ba.ID] = append(alreadyVotedMap[ba.ID], req.AgentId)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errors.New("ballot does not exist"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func doResult(w http.ResponseWriter, r *http.Request) {
	// décodage de la requête
	req, err := decodeResultRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	log.Println("Résultat pour", req.BallotId)

	// vérification que le ballot existe
	if ba, ok := ballotMap[req.BallotId]; ok {
		// vérification de la méthode de la requête
		if !checkMethod("POST", w, r) {
			return
		}

		// Tester si le vote est fini
		if !ba.DeadlineEnded() {
			w.WriteHeader(http.StatusTooEarly)
			fmt.Fprint(w, "The ballot is still in progress")
			return
		}

		// Envoyer le résultat du vote
		var res agt.ResultResponse
		var result comsoc.Alternative
		if ba.VotingFunc != "approval" {
			f := ruleMap[ba.VotingFunc]
			result, _ = f(ba.Profile)
			if ba.VotingFunc != "kemeny" {
				r := SWFMap[ba.VotingFunc]
				rank, _ := r(ba.Profile)
				for i := range rank {
					res.Ranking = append(res.Ranking, int(rank[i]))
				}
			}
		} else if ba.VotingFunc == "approval" {
			l, _ := comsoc.ApprovalSCF(ba.Profile, threeshold)
			result, _ = comsoc.TieBreak(l)
			c, _ := comsoc.ApprovalSWF(ba.Profile, threeshold)
			log.Println("Le rankinng ne marche pas pour approval")
			var rank []comsoc.Alternative
			for i := range c {
				rank = append(rank, i)
			}

			for i := range rank {
				res.Ranking = append(res.Ranking, int(rank[i]))
			}
			res.Ranking = nil
		}
		res.Winner = int(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		log.Println("  ↪   Le résultat est", res.Winner, "et le classement est", res.Ranking)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "The ballot is not found")
		return
	}
}

func LaunchServ(port string) {
	url = ":" + port
	// Attribution des règles de vote
	ruleMap["majority"] = comsoc.SCFFactory(comsoc.MajoritySWF, comsoc.TieBreak)
	ruleMap["borda"] = comsoc.SCFFactory(comsoc.BordaSWF, comsoc.TieBreak)
	ruleMap["copeland"] = comsoc.SCFFactory(comsoc.CopelandSWF, comsoc.TieBreak)
	ruleMap["kemeny"] = comsoc.FindKemenyWinner
	ruleMap["stv"] = comsoc.SCFFactory(comsoc.STV_SWF, comsoc.TieBreak)

	// Règle pour le ranking
	SWFMap["majority"] = comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreak)
	SWFMap["borda"] = comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreak)
	SWFMap["copeland"] = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreak)
	SWFMap["stv"] = comsoc.SWFFactory(comsoc.STV_SWF, comsoc.TieBreak)
	// Pas de classement pour kemeny
	// Approval est géré différemment par doResult()

	// Création du multiplexeur
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", doNewBallot)
	mux.HandleFunc("/vote", doVote)
	mux.HandleFunc("/result", doResult)

	// Démarage du serveur
	s := &http.Server{
		Addr:           url,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	log.Println("démarrage du serveur...")

	// Attente de requêtes
	go log.Fatal(s.ListenAndServe())

	// Ctrl + C pour arrêter le serveur
	fmt.Scanln()
}
