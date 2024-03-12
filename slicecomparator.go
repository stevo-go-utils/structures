package structures

import "slices"

func CompareSlices[T any](s1 []any, s2 []any) (added []any, removed []any) {
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
