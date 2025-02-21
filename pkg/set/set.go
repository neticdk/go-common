package set

import "fmt"

type Set[E comparable] map[E]struct{}

// New creates a new set with the given values.
//
// Example:
//
// s := New(1, 2, 3)
func New[E comparable](vals ...E) Set[E] {
	s := Set[E]{}
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

// Contains returns true if the set contains the given value.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Contains(2))
//	true
func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

// Add adds the given values to the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Add(4, 5)
//	fmt.Println(s.Members())
//	[1 2 3 4 5]
func (s Set[E]) Add(vals ...E) {
	for _, v := range vals {
		s[v] = struct{}{}
	}
}

// AddImmutable adds the given values to the set and returns a new set.
//
// Example:
//
// s := NewSet(1, 2, 3)
// s2 := s.AddImmutable(4, 5)
// fmt.Println(s.Members())
// [1 2 3]
// fmt.Println(s2.Members())
// [1 2 3 4 5]
func (s Set[E]) AddImmutable(vals ...E) Set[E] {
	n := New(s.Members()...)
	n.Add(vals...)
	return n
}

// Remove removes the given values from the set.
//
// Example:
//
// s := NewSet(1, 2, 3)
// s.Remove(2, 3)
// fmt.Println(s.Members())
// [1]
func (s Set[E]) Remove(vals ...E) {
	for _, v := range vals {
		delete(s, v)
	}
}

// RemoveImmutable removes the given values from the set and returns
// a new set.
//
// Example:
//
// s := NewSet(1, 2, 3)
// s2 := s.RemoveImmutable(2, 3)
// fmt.Println(s1.Members())
// [1 2 3]
// fmt.Println(s2.Members())
// [1]
func (s Set[E]) RemoveImmutable(vals ...E) Set[E] {
	n := New(s.Members()...)
	n.Remove(vals...)
	return n
}

// Clear removes all values from the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Clear()
//	fmt.Println(s.Members())
//	[]
func (s Set[E]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Members returns the members of the set as a slice.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Members())
//	[1 2 3]
func (s Set[E]) Members() []E {
	result := make([]E, 0, len(s))
	for v := range s {
		result = append(result, v)
	}
	return result
}

// String returns a string representation of the set. Order is not
// guaranteed.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s)
//	[1 2 3] *or* [3 2 1] *or* [2 1 3] *or* ...
func (s Set[E]) String() string {
	return fmt.Sprintf("%v", s.Members())
}

// Union returns the union of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Union(s2))
//	[1 2 3 4 5]
func (s Set[E]) Union(s2 Set[E]) Set[E] {
	result := New(s.Members()...)
	result.Add(s2.Members()...)
	return result
}

// Intersection returns the intersection of the set with another
// set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Intersection(s2))
//	[3]
func (s Set[E]) Intersection(s2 Set[E]) Set[E] {
	result := New[E]()
	for _, v := range s.Members() {
		if s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns the difference of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Difference(s2))
//	[1]
func (s Set[E]) Difference(s2 Set[E]) Set[E] {
	result := New(s.Members()...)
	for _, v := range s2.Members() {
		delete(result, v)
	}
	return result
}

// IsSubsetOf returns true if the set is a subset of another set.
//
// Example:
//
//	s1 := NewSet(1, 2)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.IsSubsetOf(s2))
//	true
func (s Set[E]) IsSubsetOf(s2 Set[E]) bool {
	for _, v := range s.Members() {
		if !s2.Contains(v) {
			return false
		}
	}
	return true
}

// IsSupersetOf returns true if the set is a superset of another
// set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(1, 2)
//	fmt.Println(s1.IsSupersetOf(s2))
//	true
func (s Set[E]) IsSupersetOf(s2 Set[E]) bool {
	return s2.IsSubsetOf(s)
}

// Equal returns true if the set is equal to another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.Equal(s2))
//	true
func (s Set[E]) Equal(s2 Set[E]) bool {
	// return len(s) == len(s2) && s.IsSubset(s2)
	return s.IsSubsetOf(s2) && s.IsSupersetOf(s2)
}
