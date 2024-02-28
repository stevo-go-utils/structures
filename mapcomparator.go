package structures

type MapComparator[K comparable, V any] struct {
	cur map[K]V
}

func NewMapComparator[K comparable, V any](data ...map[K]V) *MapComparator[K, V] {
	if len(data) == 1 {
		return &MapComparator[K, V]{
			cur: data[0],
		}
	}
	return &MapComparator[K, V]{
		cur: make(map[K]V),
	}
}

func (m *MapComparator[K, V]) Data() map[K]V {
	return m.cur
}

func (m *MapComparator[K, V]) CompareTo(data map[K]V) (addedKeys []K, removedKeys []K) {
	for k := range m.cur {
		if _, ok := data[k]; !ok {
			removedKeys = append(removedKeys, k)
		}
	}
	for k := range data {
		if _, ok := m.cur[k]; !ok {
			addedKeys = append(addedKeys, k)
		}
	}
	return
}

func (m *MapComparator[K, V]) CompareToSafe(data *SafeMap[K, V]) (addedKeys []K, removedKeys []K) {
	for k := range m.cur {
		if _, ok := data.Get(k); !ok {
			removedKeys = append(removedKeys, k)
		}
	}
	data.ForEach(func(k K, v V) {
		if _, ok := m.cur[k]; !ok {
			addedKeys = append(addedKeys, k)
		}
	})
	return
}

func (m *MapComparator[K, V]) RemovedKeys(data map[K]V) (removedKeys []K) {
	for k := range m.cur {
		if _, ok := data[k]; !ok {
			removedKeys = append(removedKeys, k)
		}
	}
	return
}

func (m *MapComparator[K, V]) RemovedKeysSafe(data *SafeMap[K, V]) (removedKeys []K) {
	for k := range m.cur {
		if _, ok := data.Get(k); !ok {
			removedKeys = append(removedKeys, k)
		}
	}
	return
}

func (m *MapComparator[K, V]) AddedKeys(data map[K]V) (addedKeys []K) {
	for k := range data {
		if _, ok := m.cur[k]; !ok {
			addedKeys = append(addedKeys, k)
		}
	}
	return
}

func (m *MapComparator[K, V]) AddedKeysSafe(data *SafeMap[K, V]) (addedKeys []K) {
	data.ForEach(func(k K, v V) {
		if _, ok := m.cur[k]; !ok {
			addedKeys = append(addedKeys, k)
		}
	})
	return
}

func (m *MapComparator[K, V]) Set(data map[K]V) {
	m.cur = data
}

type SafeMapComparator[K comparable, V any] struct {
	cur *SafeMap[K, V]
}

func NewSafeMapComparator[K comparable, V any](data ...*SafeMap[K, V]) *SafeMapComparator[K, V] {
	if len(data) == 1 {
		return &SafeMapComparator[K, V]{
			cur: data[0],
		}
	}
	return &SafeMapComparator[K, V]{
		cur: NewSafeMap[K, V](),
	}
}

func (m *SafeMapComparator[K, V]) Data() *SafeMap[K, V] {
	return m.cur
}

func (m *SafeMapComparator[K, V]) CompareTo(data map[K]V) (addedKeys []K, removedKeys []K) {
	m.cur.ForEach(func(k K, v V) {
		if _, ok := data[k]; !ok {
			removedKeys = append(removedKeys, k)
		}
	})
	for k := range data {
		if _, ok := m.cur.Get(k); !ok {
			addedKeys = append(addedKeys, k)
		}
	}
	return
}

func (m *SafeMapComparator[K, V]) CompareToSafe(data *SafeMap[K, V]) (addedKeys []K, removedKeys []K) {
	m.cur.ForEach(func(k K, v V) {
		if _, ok := data.Get(k); !ok {
			removedKeys = append(removedKeys, k)
		}
	})
	data.ForEach(func(k K, v V) {
		if _, ok := m.cur.Get(k); !ok {
			addedKeys = append(addedKeys, k)
		}
	})
	return
}

func (m *SafeMapComparator[K, V]) RemovedKeys(data map[K]V) (removedKeys []K) {
	m.cur.ForEach(func(k K, v V) {
		if _, ok := data[k]; !ok {
			removedKeys = append(removedKeys, k)
		}
	})
	return
}

func (m *SafeMapComparator[K, V]) RemovedKeysSafe(data *SafeMap[K, V]) (removedKeys []K) {
	m.cur.ForEach(func(k K, v V) {
		if _, ok := data.Get(k); !ok {
			removedKeys = append(removedKeys, k)
		}
	})
	return
}

func (m *SafeMapComparator[K, V]) AddedKeys(data map[K]V) (addedKeys []K) {
	for k := range data {
		if _, ok := m.cur.Get(k); !ok {
			addedKeys = append(addedKeys, k)
		}
	}
	return
}

func (m *SafeMapComparator[K, V]) AddedKeysSafe(data *SafeMap[K, V]) (addedKeys []K) {
	data.ForEach(func(k K, v V) {
		if _, ok := m.cur.Get(k); !ok {
			addedKeys = append(addedKeys, k)
		}
	})
	return
}

func (m *SafeMapComparator[K, V]) Set(data *SafeMap[K, V]) {
	m.cur = data
}
