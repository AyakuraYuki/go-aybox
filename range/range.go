package ay_range

import (
	"cmp"
)

func create[C cmp.Ordered](lower C, lowerBound Bound, upper C, upperBound Bound) Range[C] {
	return rangeImpl[C]{
		lower:      lower,
		lowerBound: lowerBound,
		upper:      upper,
		upperBound: upperBound,
	}
}

type Range[C cmp.Ordered] interface {
	HasLowerBound() bool
	LowerEndpoint() C
	LowerBoundType() Bound

	HasUpperBound() bool
	UpperEndpoint() C
	UpperBoundType() Bound

	IsEmpty() bool
	Contains(value C) bool
	ContainsAll(elems []C) bool
	Encloses(other Range[C]) bool
	IsConnected(other Range[C]) bool
	Intersection(connectedRange Range[C]) Range[C]
	Gap(other Range[C]) Range[C]
	Span(other Range[C]) Range[C]

	Equals(other Range[C]) bool

	LowerBound() (C, Bound)
	UpperBound() (C, Bound)
}

type rangeImpl[C cmp.Ordered] struct {
	lower      C
	lowerBound Bound
	upper      C
	upperBound Bound
}

func (r rangeImpl[C]) HasLowerBound() bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) LowerEndpoint() C {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) LowerBoundType() Bound {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) HasUpperBound() bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) UpperEndpoint() C {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) UpperBoundType() Bound {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Contains(value C) bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) ContainsAll(elems []C) bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Encloses(other Range[C]) bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) IsConnected(other Range[C]) bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Intersection(connectedRange Range[C]) Range[C] {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Gap(other Range[C]) Range[C] {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Span(other Range[C]) Range[C] {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) Equals(other Range[C]) bool {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) LowerBound() (C, Bound) {
	//TODO implement me
	panic("implement me")
}

func (r rangeImpl[C]) UpperBound() (C, Bound) {
	//TODO implement me
	panic("implement me")
}
