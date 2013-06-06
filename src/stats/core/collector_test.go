package core

import (
	"fmt"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	obj := &StatsObject{}
	obj.UserId = 1
	obj.Timestamp = time.Now().Unix()

	Collect(obj)

	for _, v := range _all[obj.UserId]._stats {
		fmt.Println(v)
	}
}

func BenchmarkCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 512; j++ {
			obj := &StatsObject{}
			obj.UserId = int32(i)
			obj.Timestamp = time.Now().Unix()
			Collect(obj)
		}
	}
}
