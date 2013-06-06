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
