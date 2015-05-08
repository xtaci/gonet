package bitset

type BitSet struct {
	_size uint32
	_bits []byte
}

func New(num_bits uint32) *BitSet {
	bs := &BitSet{}
	byte_len := num_bits/8 + 1
	bs._size = byte_len * 8
	bs._bits = make([]byte, byte_len)

	return bs
}

//----------------------------------------------------------  set 1 to position [bit]
func (bs *BitSet) Set(bit uint32) {
	if bit >= bs._size {
		return
	}

	n := bit / 8
	off := bit % 8

	bs._bits[n] |= 128 >> off
}

//----------------------------------------------------------  set 0 to position [bit]
func (bs *BitSet) Unset(bit uint32) {
	if bit >= bs._size {
		return
	}

	n := bit / 8
	off := bit % 8

	bs._bits[n] &= ^(128 >> off)
}

//---------------------------------------------------------- test wheather a bit is set
func (bs *BitSet) Test(bit uint32) bool {
	if bit >= bs._size {
		return false
	}

	n := bit / 8
	off := bit % 8

	if bs._bits[n]&(128>>off) != 0 {
		return true
	} else {
		return false
	}
}
