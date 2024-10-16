package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	count = make(Count)
	err = CheckProfile(p)
	if err != nil {
		return nil, err
	}
	score_max := len(p[0])
	for _, pref := range p {
		for _, alt := range pref {
			rank := Rank(alt, pref)
			count[alt] += score_max - rank
		}
	}
	return count, err
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	swf, err := BordaSWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(swf), nil
}
