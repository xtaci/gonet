package core

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkAdd(b *testing.B) {
	params := "1"
	for i := 0; i < b.N; i++ {
		Add(1, int32(i), time.Now().Unix(), []byte(params))
	}

	fmt.Println("num of events:", len(_events))

	for i := 0; i < b.N; i++ {
		Cancel(int32(i))
	}
}
