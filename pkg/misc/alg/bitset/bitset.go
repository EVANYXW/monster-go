package bitset

type BitSet struct {
	Size uint32
	Bits []byte
}

func New(numBits uint32) *BitSet {
	bs := &BitSet{}
	byte_len := numBits/8 + 1
	bs.Size = byte_len * 8
	bs.Bits = make([]byte, byte_len)

	return bs
}

//----------------------------------------------------------  set 1 to position [bit]
func (bs *BitSet) Set(bit uint32) {
	if bit >= bs.Size {
		return
	}

	n := bit / 8
	off := bit % 8

	bs.Bits[n] |= 128 >> off
}

//----------------------------------------------------------  set 0 to position [bit]
func (bs *BitSet) Unset(bit uint32) {
	if bit >= bs.Size {
		return
	}

	n := bit / 8
	off := bit % 8

	bs.Bits[n] &^= 128 >> off
}

//---------------------------------------------------------- test wheather a bit is set
func (bs *BitSet) Test(bit uint32) bool {
	if bit >= bs.Size {
		return false
	}

	n := bit / 8
	off := bit % 8

	if bs.Bits[n]&(128>>off) != 0 {
		return true
	} else {
		return false
	}
}
