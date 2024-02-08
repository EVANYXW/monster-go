package rand

var _a uint32 = 1664525
var _c uint32 = 1013904223

type Rand struct {
	X0 uint32
}

func (r *Rand) Seed(seed uint32) {
	r.X0 = seed
}

func (r *Rand) Rand() uint32 {
	r.X0 = _a*r.X0 + _c
	return r.X0
}

// [A,B]
func (r *Rand) RandAB(A, B uint32) uint32 {
	if A == B {
		return B
	}

	return r.Rand()%(B-A+1) + A
}
