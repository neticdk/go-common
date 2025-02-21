package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	s := NewSet(1, 2, 3)
	assert.Equal(t, 3, len(s))
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
}

func TestSetContains(t *testing.T) {
	s := NewSet(1, 2, 3)
	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(4))
}

func TestSetAdd(t *testing.T) {
	s := NewSet(1, 2, 3)
	s.Add(4, 5)
	assert.Equal(t, 5, len(s))
	assert.True(t, s.Contains(4))
	assert.True(t, s.Contains(5))
}

func TestSetAddImmutable(t *testing.T) {
	s := NewSet(1, 2, 3)
	s2 := s.AddImmutable(4, 5)
	assert.Equal(t, 3, len(s))
	assert.Equal(t, 5, len(s2))
	assert.True(t, s2.Contains(4))
	assert.True(t, s2.Contains(5))
}

func TestSetRemove(t *testing.T) {
	s := NewSet(1, 2, 3)
	s.Remove(2, 3)
	assert.Equal(t, 1, len(s))
	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	assert.False(t, s.Contains(3))
}

func TestSetRemoveImmutable(t *testing.T) {
	s := NewSet(1, 2, 3)
	s2 := s.RemoveImmutable(2, 3)

	assert.Equal(t, 3, len(s))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))

	assert.Equal(t, 1, len(s2))
	assert.True(t, s2.Contains(1))
	assert.False(t, s2.Contains(2))
	assert.False(t, s2.Contains(3))
}

func TestSetClear(t *testing.T) {
	s := NewSet(1, 2, 3)
	s.Clear()
	assert.Equal(t, 0, len(s))
}

func TestSetMembers(t *testing.T) {
	s := NewSet(1, 2, 3)
	members := s.Members()
	assert.ElementsMatch(t, []int{1, 2, 3}, members)

	s = NewSet[int]()
	members = s.Members()
	assert.ElementsMatch(t, []int{}, members)
}

func TestSetString(t *testing.T) {
	testCases := []struct {
		name     string
		set      Set[int]
		expected []string
	}{
		{"Three elements", NewSet(1, 2, 3), []string{"[1 2 3]", "[1 3 2]", "[2 1 3]", "[2 3 1]", "[3 1 2]", "[3 2 1]"}},
		{"No elements", NewSet[int](), []string{"[]"}},
		{"One element", NewSet(100), []string{"[100]"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var found bool
			for _, e := range tc.expected {
				if e == tc.set.String() {
					found = true
				}
			}
			assert.True(t, found)
		})
	}
}

func TestSetUnion(t *testing.T) {
	s1 := NewSet(1, 2, 3)
	s2 := NewSet(3, 4, 5)
	u := s1.Union(s2)

	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, u.Members())
}

func TestSetIntersection(t *testing.T) {
	s1 := NewSet(1, 2, 3)
	s2 := NewSet(3, 4, 5)
	i := s1.Intersection(s2)
	assert.ElementsMatch(t, []int{3}, i.Members())
}

func TestSetDifference(t *testing.T) {
	s1 := NewSet(1, 2, 3)
	s2 := NewSet(3, 4, 5)
	d := s1.Difference(s2)
	assert.ElementsMatch(t, []int{1, 2}, d.Members())
}

func TestSetIsSubsetOf(t *testing.T) {
	testCases := []struct {
		name     string
		s1       Set[int]
		s2       Set[int]
		expected bool
	}{
		{"s1 is subset", NewSet(1, 2), NewSet(1, 2, 3), true},
		{"s1 not subset", NewSet(1, 2, 4), NewSet(1, 2, 3), false},
		{"Empty set", NewSet[int](), NewSet(1, 2, 3), true},
		{"Same set", NewSet(1, 2, 3), NewSet(1, 2, 3), true},
		{"s1 is superset", NewSet(1, 2, 3), NewSet(1, 2), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.s1.IsSubsetOf(tc.s2))
		})
	}
}

func TestSetIsSupersetOf(t *testing.T) {
	s1 := NewSet(1, 2, 3)
	s2 := NewSet(1, 2)
	assert.True(t, s1.IsSupersetOf(s2))

	s3 := NewSet(1, 2, 4)
	assert.False(t, s1.IsSupersetOf(s3))
}

func TestSetEqual(t *testing.T) {
	s1 := NewSet(1, 2, 3)
	s2 := NewSet(1, 2, 3)
	assert.True(t, s1.Equal(s2))

	s3 := NewSet(1, 2)
	assert.False(t, s1.Equal(s3))
}
