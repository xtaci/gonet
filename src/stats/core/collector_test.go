package core

import (
	"fmt"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	obj := &StatsObject{}
	obj.Type = TYPE_SUM
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
			obj.Type = TYPE_SUM
			obj.Key = "样本" + fmt.Sprint(j%10)
			obj.Value = float32(j)
			obj.Timestamp = time.Now().Unix()
			Collect(int32(i), obj)
		}

		fmt.Println(_archive(int32(i), _all[int32(i)]))
	}
}
