package ipc

import (
	"encoding/json"
	"fmt"
)

import (
	. "types"
)

func IPC_ping(sess *Session, obj *IPCObject) {
	var str string
	err := json.Unmarshal(obj.Object, &str)
	if err == nil {
		if obj.Sender != -1 {
			Send(-1, obj.Sender, SERVICE_PING, str)
		} else {
			fmt.Printf("received ping value: %v\n", str)
		}
	}
}
