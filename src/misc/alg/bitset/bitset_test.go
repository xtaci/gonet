package bitset

import (
	"fmt"
	"testing"
)

func TestBitSet(t *testing.T) {
	LEN := uint32(1024)
	bs := New(LEN)
	for i := uint32(0); i < LEN; i++ {
		bs.Set(i)
	}

	fmt.Println(bs._bits)
	for i := uint32(0); i < LEN; i++ {
		if !bs.Test(i) {
			t.Fatal("bitset failed")
		}

		bs.Unset(i)
	}

	for i := uint32(0); i < LEN; i++ {
		if bs.Test(i) {
			t.Fatal("bitset failed")
		}
	}

	fmt.Println(bs._bits)
}
