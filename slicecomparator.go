package structures

import "slices"

func CompareSlices[T comparable](s1 []T, s2 []T) (added []T, removed []T) {
	for _, v1 := range s1 {
		if !slices.Contains(s2, v1) {
			removed = append(removed, v1)
		}
	}
	for _, v2 := range s2 {
		if !slices.Contains(s1, v2) {
			added = append(added, v2)
		}
	}
	return
}
