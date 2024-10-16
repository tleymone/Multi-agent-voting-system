package comsoc

import "errors"

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	if len(thresholds) != len(p) {
		return nil, errors.New("thresholds table has not the same length as the profile's length")
	}
	count = make(Count)
	err = CheckProfile(p)
	if err != nil {
		return nil, err
	}
	for i, pref := range p {
		for _, alt := range pref[0:thresholds[i]] {
			count[alt] += 1
		}
	}
	return count, err
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	swf, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return MaxCount(swf), nil
}
