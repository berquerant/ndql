package jsonmap

import (
	"encoding/json"
	"strings"
)

type Map[T any] struct {
	value    T
	hasValue bool
	d        map[string]*Map[T]
}

func (m *Map[T]) MarshalJSON() ([]byte, error) {
	v := map[string]any{
		"hasValue": m.hasValue,
		"d":        m.d,
	}
	if m.hasValue {
		v["value"] = m.value
	}
	return json.Marshal(v)
}

func (m *Map[T]) UnmarshalJSON(b []byte) error {
	type X struct {
		Value    T                  `json:"value,omitempty"`
		HasValue bool               `json:"hasValue"`
		D        map[string]*Map[T] `json:"d"`
	}
	var x X
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*m = Map[T]{
		value:    x.Value,
		hasValue: x.HasValue,
		d:        x.D,
	}
	return nil
}

func newMap[T any]() *Map[T] {
	return &Map[T]{
		d: map[string]*Map[T]{},
	}
}

func (m *Map[T]) Value() (T, bool)             { return m.value, m.hasValue }
func (m *Map[T]) Children() map[string]*Map[T] { return m.d }

func (m *Map[T]) Get(key string) (*Map[T], bool) {
	return m.get(strings.Split(key, ".")...)
}

func (m *Map[T]) get(key ...string) (*Map[T], bool) {
	if len(key) == 0 {
		return nil, false
	}
	d, ok := m.d[key[0]]
	if !ok {
		return nil, false
	}
	if len(key) == 1 {
		return d, true
	}
	return d.get(key[1:]...)
}

func (m *Map[T]) set(value T, key ...string) {
	if len(key) == 0 {
		return
	}
	k := key[0]
	if len(key) == 1 {
		if d, ok := m.d[k]; ok {
			d.value = value
		} else {
			m.d[k] = newMap[T]()
			m.d[k].value = value
		}
		m.d[k].hasValue = true
		return
	}

	if d, ok := m.d[k]; ok {
		d.set(value, key[1:]...)
		return
	}
	m.d[k] = newMap[T]()
	m.d[k].set(value, key[1:]...)
}

func (m *Map[T]) Set(key string, value T) {
	m.set(value, strings.Split(key, ".")...)
}

func NewMap[T any]() *Map[T] {
	return newMap[T]()
}
