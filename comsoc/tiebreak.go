package comsoc

import (
	"errors"
)

func TieBreak(alts []Alternative) (Alternative, error) {
	if len(alts) == 0 {
		return -1, errors.New("there are no alternatives to tie-break")
	}
	return alts[0], nil
}

func TieBreakFactory(alts []Alternative) func([]Alternative) (Alternative, error) {

	f := func(alts2 []Alternative) (Alternative, error) {
		if len(alts) == 0 {
			return -1, errors.New("there are no alternatives to tie-break")
		}
		var fav Alternative = -1
		for i, alt1 := range alts2 {
			for j := i + 1; j < len(alts2); j++ {
				if IsPref(alt1, fav, alts) {
					if IsPref(alt1, alts2[j], alts) {
						fav = alt1
					} else {
						fav = alts2[j]
					}
				}
			}
		}
		return fav, nil
	}
	return f
}
