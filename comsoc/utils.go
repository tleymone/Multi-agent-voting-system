package comsoc

import "errors"

func Rank(alt Alternative, prefs []Alternative) int {
	for i := 0; i < len(prefs); i++ {
		if prefs[i] == alt {
			return i
		}
	}
	return -1
}

func IsPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	i1 := Rank(alt1, prefs)
	i2 := Rank(alt2, prefs)

	if (i1 != -1 && i2 != -1 && i1 < i2) || (i2 == -1) {
		return true
	} else {
		return false
	}
}

func MaxCount(count Count) (bestAlts []Alternative) {
	max := 0
	for _, v := range count {
		if v > max {
			max = v
		}
	}
	for i, v := range count {
		if v == max {
			bestAlts = append(bestAlts, i)
		}
	}
	return bestAlts
}

func CheckProfile(prefs Profile) error {
	len_pref := len(prefs[0])
	for _, pref := range prefs {
		if len(pref) != len_pref {
			return errors.New("all prefs are not complete")
		}
	}
	for _, pref := range prefs {
		for i, alt := range pref {
			for j := i + 1; j < len(pref); j++ {
				if pref[j] == alt {
					return errors.New("an alternative appears several times in at least one preference")
				}
			}
		}
	}
	return nil
}

func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	err := CheckProfile(prefs)
	if err != nil {
		return err
	}
	for _, alt := range alts {
		for _, pref := range prefs {
			b := false
			if len(pref) != len(alts) {
				return errors.New("every alternative must appear exactly once in a preference")
			}
			for _, v := range pref {
				if v == alt {
					b = true
				}
			}

			if !b {
				return errors.New("unknown alternative in one of the preferences")
			}

		}
	}
	return nil
}
