package internal

// A go map that keeps its order intact
type OrderedMap[V any] struct {
	keys   []string
	values map[string]V
}

func NewOrderedMap[V any]() *OrderedMap[V] {
	return &OrderedMap[V]{
		keys:   []string{},
		values: map[string]V{},
	}
}

func (omap *OrderedMap[V]) IsEmpty() bool {
	return len(omap.keys) == 0
}

func (omap *OrderedMap[V]) Add(key string, value V) {
	omap.keys = append(omap.keys, key)
	omap.values[key] = value
}

func (omap *OrderedMap[V]) ForEach(fn func(string, V)) {
	for _, key := range omap.keys {
		fn(key, omap.values[key])
	}
}
