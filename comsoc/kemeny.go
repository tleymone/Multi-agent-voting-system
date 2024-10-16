package comsoc

func createAllPairs(pref []Alternative) (pairs [][]Alternative) {
	pairs = [][]Alternative{}
	for i, alt := range pref {
		for j := i + 1; j < len(pref); j++ {
			pairs = append(pairs, []Alternative{alt, pref[j]})
		}
	}
	return
}

func CalcDistRangement(pref1, pref2 []Alternative) (dist int) {
	//distance tau de Kendall
	dist = 0
	pairs1 := createAllPairs(pref1)
	pairs2 := createAllPairs(pref2)

	//recherche de toutes les paires différentes entre les deux listes de paires
	for _, pair1 := range pairs1 {
		for _, pair2 := range pairs2 {
			if pair1[0] == pair2[1] && pair1[1] == pair2[0] {
				//les deux paires sont dans des ordres différents, il faut donc incrémenter la distance
				dist++
			}
		}
	}
	return
}

func CalcDistProfile(pref []Alternative, p Profile) (dist int, err error) {
	err = CheckProfile(p)
	if err != nil {
		return -1, err
	}

	dist = 0
	for _, pref2 := range p {
		dist += CalcDistRangement(pref, pref2)
	}
	return
}

func FindKemenyWinner(p Profile) (Alternative, error) {
	minDist := 100000
	var minDistPref []Alternative
	for _, pref := range p {
		//pour chaque préférence, calculer sa distance avec l'ensemble du profil
		dist, err := CalcDistProfile(pref, p)
		if err != nil {
			return -1, err
		}
		//garder en mémoire la préférence avec la plus petite distance
		if dist < minDist {
			minDist = dist
			minDistPref = pref
		}
	}
	//retourner le premier candidat présent dans la préférence qui possède la plus petite distance
	return minDistPref[0], nil
}
