package forward_tbl

import (
	"encoding/json"
	"fmt"
	"testing"
	. "types"
)

func TestForward(t *testing.T) {
	obj := &IPCObject{Sender: 1}
	json_obj, _ := json.Marshal(obj)
	Push(1, json_obj)
	objs := PopAll(1)

	for k := range objs {
		fmt.Println(k, string(objs[k]))
	}
}
