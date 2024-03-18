package structures

import "slices"

func CompareMaps[K comparable, V any](m1 map[K]V, m2 map[K]V) (added []K, removed []K) {
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			removed = append(removed, k)
		}
	}
	for k := range m2 {
		if _, ok := m1[k]; !ok {
			added = append(added, k)
		}
	}
	return
}

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

func CompareSliceToMap[K comparable, V any](s1 []K, m2 map[K]V) (added []K, removed []K) {
	for _, v := range s1 {
		if _, ok := m2[v]; !ok {
			removed = append(removed, v)
		}
	}
	for k := range m2 {
		if !slices.Contains(s1, k) {
			added = append(added, k)
		}
	}
	return
}

func CompareMapToSlice[K comparable, V any](m1 map[K]V, s2 []K) (added []K, removed []K) {
	for k := range m1 {
		if !slices.Contains(s2, k) {
			removed = append(removed, k)
		}
	}
	for _, v := range s2 {
		if _, ok := m1[v]; !ok {
			added = append(added, v)
		}
	}
	return
}
