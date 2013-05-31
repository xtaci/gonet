package main

import (
	"fmt"
	"log"
)

import (
	"agent/client_protos"
	"agent/ipc"
	"misc/packet"
	. "misc/stack"
	. "types"
)

//----------------------------------------------- client protocol handle proxy
func UserRequestProxy(sess *Session, p []byte) []byte {
	defer PrintPanicStack()

	reader := packet.Reader(p)
	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	log.Printf("code:%v,user:%v\n", b, sess.User.Id)

	handle := protos.ProtoHandler[b]
	if handle != nil {
		ret, err := handle(sess, reader)
		fmt.Println(ret)
		if err == nil {
			return ret
		}
	}

	return nil
}

//----------------------------------------------- IPC proxy
func IPCRequestProxy(sess *Session, p *IPCObject) {
	defer PrintPanicStack()
	handle := ipc.IPCHandler[p.Service]
	log.Printf("ipc:%v,user:%v\n", p.Service, sess.User.Id)

	if handle != nil {
		handle(sess, p)
	}
}
