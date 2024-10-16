package comsoc

import (
	"sort"
)

func SWFFactory(C func(Profile) (Count, error), T func([]Alternative) (Alternative, error)) func(p Profile) (a []Alternative, e error) {
	f := func(p Profile) ([]Alternative, error) {
		c, e := C(p)
		if e != nil {
			return nil, e
		}
		var arr []Alternative
		for i := range c {
			arr = append(arr, i)
		}
		sort.SliceStable(arr, func(i, j int) bool {
			return c[arr[i]] > c[arr[j]]
		})
		return arr, e
	}
	return f
}

func SCFFactory(C func(Profile) (Count, error), T func([]Alternative) (Alternative, error)) func(p Profile) (a Alternative, e error) {
	f := func(p Profile) (Alternative, error) {
		c, e := C(p)
		if e != nil {
			return -1, e
		}
		max := MaxCount(c)
		a, e := T(max)
		return a, e
	}
	return f
}
