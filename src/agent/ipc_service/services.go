package ipc_service

import (
	"log"
	"runtime"
)

import (
	. "agent/ipc"
	. "types"
)

var IPCHandler map[int16]func(*Session, *IPCObject) []byte = map[int16]func(*Session, *IPCObject) []byte{
	SVC_PING: IPC_ping,
	SVC_CHAT: IPC_chat,
	SVC_KICK: IPC_kick,

	SYS_BROADCAST: SYS_broadcast,
	SYS_MULTICAST: SYS_multicast,
}

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in ipc_service module")
	}
}
