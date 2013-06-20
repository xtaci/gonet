package forward_tbl

import (
	"fmt"
	"testing"
	. "types"
)

func TestForward(t *testing.T) {
	obj := &IPCObject{SrcID: 1, DestID: 2}
	fmt.Println(Push(obj))
	objs := PopAll(2)

	for k := range objs {
		fmt.Println("retrieved:", objs[k])
	}
}
