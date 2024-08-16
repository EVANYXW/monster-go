package bloomfilter

import (
	"crypto/rand"
)

type UHash struct {
	A   int64
	B   int64
	P   int64
	MAX int64
}

func (h *UHash) Hash(X uint32) uint32 {
	hash := ((h.A*int64(X) + h.B) % h.P) % h.MAX
	return uint32(hash)
}

//---------------------------------------------------------- init a hash function
func UHashInit(numbins uint32) UHash {
	h := UHash{}
	prime, _ := rand.Prime(rand.Reader, _msb(numbins)+1)
	h.P = prime.Int64()
	h.MAX = int64(numbins)

	a, _ := rand.Int(rand.Reader, prime)
	b, _ := rand.Int(rand.Reader, prime)

	h.A = a.Int64()
	h.B = b.Int64()

	for h.A == 0 { // make sure h.A is not zero
		a, _ = rand.Int(rand.Reader, prime)
		h.A = a.Int64()
	}

	return h
}

func _msb(n uint32) int {
	r := 0
	for ; n != 0; n >>= 1 {
		r++
	}

	return r
}
