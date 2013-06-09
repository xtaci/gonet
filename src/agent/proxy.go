package main

import (
	"fmt"
	"log"
	"time"
)

import (
	"agent/client_protos"
	"agent/ipc"
	. "helper"
	"misc/packet"
	. "types"
)

//----------------------------------------------- client protocol handle proxy
func UserRequestProxy(sess *Session, p []byte) []byte {
	defer PrintPanicStack()
	now := time.Now()

	reader := packet.Reader(p)
	// read client_elapsed in milli-second since startup time-sync
	client_elapsed, err := reader.ReadU32()
	if err != nil {
		log.Println("Read timestamp failed.", err)
		return nil
	}

	if sess.LoggedIn {
		server_elapsed := now.Sub(sess.ConnectTime).Nanoseconds() / 1000
		diff := int16(server_elapsed - int64(client_elapsed))
		sess.User.LatencySamples.Add(diff)
	}

	// read protocol id
	b, err := reader.ReadU16()
	if err != nil {
		log.Println("read protocol error")
	}

	fmt.Printf("code:%v,user:%v\n", b, sess.User.Id)

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
