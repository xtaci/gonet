package forward_tbl

import (
	"fmt"
	"testing"
	. "types"
)

func TestForward(t *testing.T) {
	obj := &IPCObject{SrcID: 0, DestID: 1}
	fmt.Println(Push(obj))
	fmt.Println(Push(obj))
	objs := PopAll(1)

	for k := range objs {
		fmt.Println("retrieved:", objs[k])
	}

	if len(objs) != 2 {
		t.Fatal("forward db failed getall")
	}
}
