package main

import (
	"agent/ipc"
	"testing"
	. "types"
)

func TestIPC(t *testing.T) {
	// fake 2 user
	var sess1 Session
	sess1.User.Id = 1
	sess1.MQ = make(chan IPCObject, 10)
	var sess2 Session
	sess2.User.Id = 2
	sess2.MQ = make(chan IPCObject, 10)

	ipc.RegisterOnline(&sess1, sess1.User.Id)
	ipc.RegisterOnline(&sess2, sess2.User.Id)

	ipc.Send(1, 2, ipc.SERVICE_PING, "ABC")

	obj := <-sess2.MQ
	IPCRequestProxy(&sess2, &obj)
	obj = <-sess1.MQ
	IPCRequestProxy(&sess1, &obj)
}
