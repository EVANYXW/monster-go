package bloomfilter

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	filter := NewFilter(1000000, 100000)

	strs := []string{
		"purple_lover_04@msn.com",
		"hollie.smith@yahoo.com",
		"hollister.susan @gmail.com",
		"alexptre@gmail.com",
		"cfo@sjtu-edp.cn",
		"abms@n23h22.rev.sprintdatacenter.pl",
		"admin@facebook.com",
		"xtaci@163.com",
	}

	for k := range strs {
		filter.Set(strs[k])
	}

	if !filter.Test("xtaci@163.com") || !filter.Test("abms@n23h22.rev.sprintdatacenter.pl") || filter.Test("xxx") {
		t.Fatal("bloom filter failed")
	}
}

func BenchmarkBF(b *testing.B) {
	filter := NewFilter(1000000, 100000)

	for i := 0; i < b.N; i++ {
		filter.Test("abc")
	}
}
