package main

import (
	"log"
	"net"
	"sync/atomic"
)

import (
	"hub/core"
	"hub/protos"
	"misc/packet"
	. "types"
)

//---------------------------------------------------------- 连入GS的ID生成器
var _host_genid int32

func init() {
	log.SetPrefix("[HUB] ")
}

//---------------------------------------------------------- HUB SERVICE
func StartAgent(incoming chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	log.Printf("game server [id:%v] connected\n", hostid)

	// forward queue (IPCObject)
	forward := make(chan IPCObject, 100000)
	protos.AddServer(hostid, forward)

	// closing
	defer func() {
		protos.RemoveServer(hostid)
		core.LogoutServer(hostid)
		close(forward)

		log.Printf("game server [id:%v] disconnected\n", hostid)
	}()

	for {
		select {
		case msg, ok := <-incoming: // request from game server
			if !ok {
				return
			}

			// read seqid
			reader := packet.Reader(msg)
			seqid, err := reader.ReadU32()
			if err != nil {
				log.Println("read SEQID failed.", err)
				return
			}

			// handle request
			ret := GSProxy(hostid, reader)
			// send result
			if len(ret) != 0 {
				_send(seqid, ret, conn)
			}
		case obj := <-forward: // forwarding packets(ie. seqid == 0)
			_send(0, obj.Json(), conn)
		}
	}
}

//---------------------------------------------------------- send response
func _send(seqid uint32, data []byte, conn net.Conn) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 4)
	writer.WriteU32(seqid) // piggyback seq id
	writer.WriteRawBytes(data)

	n, err := conn.Write(writer.Data()) // write operation is assumed to be atomic
	if err != nil {
		log.Println("Error send reply to GS, bytes:", n, "reason:", err)
	}
}
