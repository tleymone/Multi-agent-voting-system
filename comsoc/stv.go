package comsoc

func remove(slice []Alternative, s int) []Alternative {
	return append(slice[:s], slice[s+1:]...)
}

//Ce STV renvoie un Count avec des scores pondérés par le nombre de tours que le candidat à passer sans se faire éliminer
//Plus un score est élevé, plus le candidat est allé loin dans le vote

func STV_SWF(p Profile) (count Count, err error) {

	//copier le profil pour ne pas modifier le profil en entrée
	var p2 = make(Profile, len(p))
	for i := range p {
		p2[i] = make([]Alternative, len(p[i]))
		copy(p2[i], p[i])
	}

	err = CheckProfile(p2)
	if err != nil {
		return nil, err
	}

	//initialiser la map à retourner
	count = make(Count)
	alt_list := p2[0]
	for _, alt := range alt_list {
		count[alt] = 0
	}

	//début des itérations
	n_alts := len(p2[0])
	for i := 0; i < n_alts; i++ {
		alt_list := p2[0]
		//initialiser une map temporaire avec tous les candidats restants dans cette itération
		var temp_count = make(Count)
		for _, alt := range alt_list {
			temp_count[alt] = 0
		}

		//compter le nombre de votes pour une itération i
		for _, pref := range p2 {
			fav := pref[0]
			temp_count[fav] += 1
		}

		//rechercher le candidat avec le moins de votes
		var worst_alt Alternative = 1000
		min_value := 10000 //valeur arbitraire d'initialisation
		for key, value := range temp_count {
			if value < min_value {
				min_value = value
				worst_alt = key
			}
		}

		//supprimer le candidat de toutes les préférences
		for i, pref := range p2 {
			s := Rank(worst_alt, pref)
			p2[i] = remove(pref, s)
		}

		//à la fin, pondérer le score le plus faible par i (numéro d'itération)
		//ainsi, le score le plus faible après toutes les itérations correspond au candidat éliminé en premier,
		//le deuxième score le plus faible, au deuxième pire candidat, etc...
		count[worst_alt] = i * temp_count[worst_alt]
	}
	return count, nil
}

func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
	swf, err := STV_SWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(swf), nil
}
