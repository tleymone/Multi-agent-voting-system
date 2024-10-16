package comsoc

func findCondorcetWinner(winner_list []Alternative, alts_list []Alternative) []Alternative {
	for _, alt := range alts_list {
		count := 0
		for _, winner := range winner_list {
			if winner == alt {
				count += 1
			}
		}
		if count == len(alts_list)-1 {
			//un gagnant de Condorcet est un candidat qui apparaît n-1 fois dans la liste des gagants (il gagne contre tout le monde sauf lui)
			return []Alternative{alt}
		}
	}
	return []Alternative{}
}

func CondorcetWinner(p Profile) ([]Alternative, error) {
	err := CheckProfile(p)
	if err != nil {
		return nil, err
	}

	alt_list := p[0]
	var winner_list []Alternative

	for i, alt1 := range alt_list {
		for j := i + 1; j < len(alt_list); j++ {
			alt2 := alt_list[j]
			n_votes_alt1 := 0
			n_votes_alt2 := 0
			for _, pref := range p {
				if Rank(alt1, pref) < Rank(alt2, pref) {
					n_votes_alt1 += 1
				} else if Rank(alt2, pref) < Rank(alt1, pref) {
					n_votes_alt2 += 1
				}
			}
			if n_votes_alt1 < n_votes_alt2 {
				winner_list = append(winner_list, alt2)
			} else if n_votes_alt2 < n_votes_alt1 {
				winner_list = append(winner_list, alt1)
			} else {
				//cas où il y a une égalité entre deux candidats, dans ce cas là il n'y a aucun gagnant de Condorcet
				return []Alternative{}, nil
			}
		}
	}
	return findCondorcetWinner(winner_list, alt_list), nil
}
