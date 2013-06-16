package AI

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkSP(b *testing.B) {
	t0 := time.Now().Unix()

	for i := int64(0); i < int64(b.N); i++ {
		SP(0, 10000, t0-i, 100)
	}
}

func BenchmarkDice(b *testing.B) {
	count := 0
	for i := 0; i < b.N; i++ {
		if Dice(0.01) {
			count++
		}
	}

	fmt.Printf("dice prob:0.01, count: %v, total: %v\n", count, b.N)
}

/*
func TestLCG(t *testing.T) {
	for i:=0;i<100;i++ {
		fmt.Println(LCG())
	}
}
*/
