package comsoc

func CopelandSWF(p Profile) (c Count, err error) {
	// test 2 à 2 pour chaque alternative
	err = CheckProfile(p)
	if err != nil {
		return nil, err
	}
	compteur2 := make(Count)

	for _, alt1 := range p[0] {
		for _, alt2 := range p[0] {
			if alt1 != alt2 {
				compteur := make(Count)
				for _, i := range p {
					if IsPref(alt1, alt2, i) {
						compteur[alt1] += 1
					} else {
						compteur[alt1] -= 1
					}
				}
				if compteur[alt1] > 0 {
					compteur2[alt1] += 1
				} else if compteur[alt1] < 0 {
					compteur2[alt1] -= 1
				}
			}
		}
	}
	return compteur2, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	swf, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(swf), nil
}
