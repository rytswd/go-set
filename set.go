// Package set defines various methods for a set.
package set

import (
	"fmt"
	"sync"
)

// Set is a set of comparables.
type Set[V comparable] struct {
	m map[V]struct{}

	mux sync.RWMutex
}

// New returns a Set from the given values.
func New[V comparable](v ...V) *Set[V] {
	s := &Set[V]{
		m:   make(map[V]struct{}),
		mux: sync.RWMutex{},
	}

	s.Insert(v...)

	return s
}

// Clone returns a new Set that a copy of `s`.
func (s *Set[V]) Clone() *Set[V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	t := New[V]()

	t.Insert(s.Values()...)

	return t
}

// Delete removes the given values from `s`.
func (s *Set[V]) Delete(v ...V) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, x := range v {
		delete(s.m, x)
	}
}

// Difference returns a Set whose values are in `s` and not in `t`.
//
// For example:
//
//	s = {a1, a2, a3}
//	t = {a1, a2, a4, a5}
//	s.Difference(t) = {a3}
//	t.Difference(s) = {a4, a5}
func (s *Set[V]) Difference(t *Set[V]) *Set[V] {
	s.mux.RLock()
	t.mux.RLock()
	defer s.mux.RUnlock()
	defer t.mux.RUnlock()

	u := New[V]()

	for k := range s.m {
		if !t.Has(k) {
			u.Insert(k)
		}
	}

	return u
}

// Intersection returns a new Set whose values are included in both `s` and `t`.
//
// For example:
//
//	s = {a1, a2}
//	t = {a2, a3}
//	s.Intersection(t) = {a2}
func (s *Set[V]) Intersection(t *Set[V]) *Set[V] {
	s.mux.RLock()
	t.mux.RLock()
	defer s.mux.RUnlock()
	defer t.mux.RUnlock()

	u := New[V]()

	var walk, other *Set[V]

	if s.Len() < t.Len() {
		walk = s
		other = t
	} else {
		walk = t
		other = s
	}

	for k := range walk.m {
		if other.Has(k) {
			u.Insert(k)
		}
	}

	return u
}

// Equal returns true iff `s` is equal to `t`.
//
// Two sets are equal if their underlying values are identical not considering
// order.
func (s *Set[V]) Equal(t *Set[V]) bool {
	s.mux.RLock()
	t.mux.RLock()
	defer s.mux.RUnlock()
	defer t.mux.RUnlock()

	return len(s.m) == len(t.m) && s.IsSuperset(t)
}

// Has returns true iff `s` contains a given value.
func (s *Set[V]) Has(v V) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	_, ok := s.m[v]
	return ok
}

// HasAny returns true iff `s` contains all the given values.
func (s *Set[V]) HasAll(v ...V) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, x := range v {
		if !s.Has(x) {
			return false
		}
	}

	return true
}

// HasAll returns true iff `s` contains any of the given values.
func (s *Set[V]) HasAny(v ...V) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, x := range v {
		if s.Has(x) {
			return true
		}
	}

	return false
}

// Insert adds the given values to `s`.
func (s *Set[V]) Insert(v ...V) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, x := range v {
		s.m[x] = struct{}{}
	}
}

// IsSuperset returns true iff `t` is a superset of `s`.
func (s *Set[V]) IsSuperset(t *Set[V]) bool {
	s.mux.RLock()
	t.mux.RLock()
	defer s.mux.RUnlock()
	defer t.mux.RUnlock()

	for k := range t.m {
		if !s.Has(k) {
			return false
		}
	}

	return true
}

// Len returns the size of `s`.
func (s *Set[V]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return len(s.m)
}

// Pop returns a single value randomly chosen and removes it from `s`.
func (s *Set[V]) PopAny() (v V, _ bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for k := range s.m {
		delete(s.m, k)
		return k, true
	}

	return v, false
}

// String implements fmt.Stringer.
func (s *Set[V]) String() string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return fmt.Sprint(s.Values())
}

// Values returns the underlying values of `s`.
func (s *Set[V]) Values() []V {
	s.mux.RLock()
	defer s.mux.RUnlock()

	v := make([]V, 0, len(s.m))

	for k := range s.m {
		v = append(v, k)
	}

	return v
}

// Union returns a new Set whose values are included in either `s` or `t`.
//
// For example:
//
//	s = {a1, a2}
//	t = {a3, a4}
//	s.Union(t) = {a1, a2, a3, a4}
//	t.Union(s) = {a1, a2, a3, a4}
func (s *Set[V]) Union(t *Set[V]) *Set[V] {
	u := s.Clone()

	s.mux.RLock()
	t.mux.RLock()
	defer s.mux.RUnlock()
	defer t.mux.RUnlock()

	for k := range t.m {
		u.Insert(k)
	}

	return u
}