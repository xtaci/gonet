package core

import (
	"fmt"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	obj := &StatsObject{}
	obj.Timestamp = time.Now().Unix()

	userid := int32(1)
	Collect(userid, obj)

	for _, v := range _all[userid]._stats {
		fmt.Println(v)
	}
}

func BenchmarkCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 200; j++ {
			obj := &StatsObject{}
			obj.Timestamp = time.Now().Unix()
			Collect(int32(i), obj)
		}
	}

}
