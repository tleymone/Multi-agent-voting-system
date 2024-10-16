package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"td3/agt"
	"td3/agt/voteragent"
	"td3/comsoc"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Lancement du serveur
	port := "8080"

	// Création du ballot
	nbVot := 20                                          // nombres des votants
	votList := []string{}                                // liste des votants
	voterAgts := make([]voteragent.VoterAgent, 0, nbVot) // liste des agents
	for i := 1; i <= nbVot; i++ {
		votList = append(votList, fmt.Sprintf("ag_id%01d", i))
	} // liste des votants du ballot
	deadline := time.Now().Add(5 * time.Second) // deadline fixée à 5 secondes

	// création du JSON pour la requête
	ballot := agt.NewBallotRequest{
		Rule:     "borda",
		Deadline: deadline.Format(time.RFC3339),
		VoterIds: votList,
		NbAlts:   5,
	}

	json_data, err := json.Marshal(ballot)
	if err != nil {
		log.Fatal(err)
	}
	// Envoi de la requête
	resp, err := http.Post("http://localhost:"+port+"/new_ballot", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err, resp)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println("Le ballot", res["ballot-id"], "a été créé, la règle est", ballot.Rule)

	options := []int{}

	// Création des votants
	for i := 0; i < nbVot; i++ {
		name := fmt.Sprintf("ag_id%0d", i)
		// créer une préférence aléatoire
		prefs := make([]comsoc.Alternative, ballot.NbAlts)
		perm := rand.Perm(ballot.NbAlts)
		for i, v := range perm {
			prefs[v] = comsoc.Alternative(i) + 1
		}
		agt := voteragent.NewVoterAgent(name, port, prefs, options)
		voterAgts = append(voterAgts, *agt)
	}

	for _, agt := range voterAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt voteragent.VoterAgent) {
			go agt.Start(fmt.Sprint(res["ballot-id"]))
		}(agt)
	}

	// Pause pour avoir le temps de voter
	time.Sleep(time.Until(deadline))
	log.Println("La deadline est passé !")
	// Demande des résultats

	// création du JSON pour la requête
	result := agt.ResultRequest{
		BallotId: fmt.Sprint(res["ballot-id"]),
	}

	json_data, err = json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = http.Post("http://localhost:"+port+"/result", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println("Le gagnant est : ", res["winner"], "et le classement est", res["ranking"])
}
