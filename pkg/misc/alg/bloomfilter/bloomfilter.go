package bloomfilter

import (
	"github.com/evanyxw/monster-go/pkg/misc/alg/bitset"
	"github.com/evanyxw/monster-go/pkg/misc/naming"
)

const HASH_FUNCS = 8

type BloomFilter struct {
	NBits  uint32
	NProb  uint32
	Hashes []UHash
	BitSet *bitset.BitSet
}

func NewFilter(M, N uint32) *BloomFilter {
	filter := &BloomFilter{}
	filter.NBits = M
	filter.NProb = N
	filter.Hashes = make([]UHash, HASH_FUNCS)
	filter.BitSet = bitset.New(M)

	for i := 0; i < HASH_FUNCS; i++ {
		filter.Hashes[i] = UHashInit(M)
	}
	return filter
}

func (filter *BloomFilter) Set(str string) {
	v := naming.FNV1a(str)
	for i := 0; i < HASH_FUNCS; i++ {
		hash := filter.Hashes[i].Hash(v)
		filter.BitSet.Set(hash)
	}
}

func (filter *BloomFilter) Test(str string) bool {
	v := naming.FNV1a(str)
	for i := 0; i < HASH_FUNCS; i++ {
		hash := filter.Hashes[i].Hash(v)
		if !filter.BitSet.Test(hash) {
			return false
		}
	}

	return true
}
