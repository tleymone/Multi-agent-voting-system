package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	count = make(Count)
	err = CheckProfile(p)
	if err != nil {
		return nil, err
	}
	for _, pref := range p {
		fav := pref[0]
		count[fav] += 1
	}
	return count, err
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	swf, err := MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(swf), nil
}
