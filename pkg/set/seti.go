package set

// SetI is the interface used by all the Set types.
type SetI[E comparable] interface {
	Contains(E) bool
	ContainsAll(...E) bool
	Add(...E)
	AddImmutable(...E) SetI[E]
	Remove(...E)
	Clear()
	Members() []E
	String() string
	Union(SetI[E]) SetI[E]
	Intersection(SetI[E]) SetI[E]
	Difference(SetI[E]) SetI[E]
	IsSubsetOf(SetI[E]) bool
	IsSupersetOf(SetI[E]) bool
	Equal(SetI[E]) bool
	Len() int
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
