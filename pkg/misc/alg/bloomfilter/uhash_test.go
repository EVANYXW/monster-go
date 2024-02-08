package bloomfilter

import (
	"fmt"
	"testing"
	"time"
)

func TestUHash(t *testing.T) {
	H1 := UHashInit(1000000)
	H2 := UHashInit(1000000)
	num := time.Now().Unix()
	fmt.Println(H1.Hash(uint32(num)))
	fmt.Println(H2.Hash(uint32(num)))
	fmt.Println("h1: prime", H1.P, "a", H1.A, "b", H1.B)
	fmt.Println("h2: prime", H2.P, "a", H2.A, "b", H2.B)
	fmt.Println("msb", _msb(1000000))
}
