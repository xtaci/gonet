package forward_tbl

import (
	"fmt"
	"testing"
	. "types"
)

func TestForward(t *testing.T) {
	obj := &IPCObject{Sender: 1}
	Push(1, obj.Json())
	objs := PopAll(1)

	for k := range objs {
		fmt.Println(k, string(objs[k]))
	}
}
