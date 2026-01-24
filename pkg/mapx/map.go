package mapx

import (
	"encoding/json"
	"maps"
	"slices"
	"sync"
)

// Wrapped map[K]V; safe for concurrent use.
type Map[K comparable, V any] struct {
	mux sync.RWMutex
	m   map[K]V
}

func NewMap[K comparable, V any](m map[K]V) *Map[K, V] {
	if m == nil {
		m = map[K]V{}
	}
	return &Map[K, V]{
		m: m,
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.m[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.m, key)
}

func (m *Map[K, V]) Clone() *Map[K, V] {
	m.mux.Lock()
	defer m.mux.Unlock()
	return &Map[K, V]{
		m: maps.Clone(m.m),
	}
}

func (m *Map[K, V]) Merge(other *Map[K, V]) {
	m.mux.Lock()
	other.mux.RLock()
	defer func() {
		other.mux.RUnlock()
		m.mux.Unlock()
	}()
	for k, v := range other.m {
		m.m[k] = v
	}
}

// Unwrap returns a cloned internal map.
func (m *Map[K, V]) Unwrap() map[K]V { return m.Clone().m }

func (m *Map[K, V]) Keys() []K {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return slices.Collect(maps.Keys(m.m))
}

func (m *Map[K, V]) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return len(m.m)
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return json.Marshal(m.m)
}

func (m *Map[K, V]) UnmarshalJSON(b []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	var d map[K]V
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	m.m = d
	return nil
}
