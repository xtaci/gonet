package helper

import (
	"testing"
)

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		println(<-LCG)
	}
}
