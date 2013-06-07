package main

import (
	"log"
	"net"
	"sync/atomic"
)

import (
	. "helper"
	"hub/protos"
	"misc/packet"
)

const (
	MAXCHAN = 100000
)

//---------------------------------------------------------- 连入主机的ID标示
var _host_genid int32

func init() {
	log.SetPrefix("[HUB]")
}

//---------------------------------------------------------- Hub processing
func HubAgent(incoming chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	// forward buffer
	forward := make(chan []byte, MAXCHAN)
	// output buffer
	output := make(chan []byte, MAXCHAN)

	protos.ServerLock.Lock()
	protos.Servers[hostid] = forward // 转发消息队列
	protos.ServerLock.Unlock()

	log.Printf("server id:%v connected\n", hostid)

	go _write_routine(output, conn)

	defer func() {
		protos.ServerLock.Lock()
		delete(protos.Servers, hostid)
		protos.ServerLock.Unlock()

		close(forward)
		close(output)

		log.Printf("server id:%v disconnected\n", hostid)
	}()

	// HUB只处理两类消息，来自GS的，和转发来自其他GS的IPC消息
	// 转发IPC消息时,seqid为0,GS检查如果seqid为0，则为forward消息
	for {
		select {
		case msg, ok := <-incoming:
			if !ok {
				return
			}

			reader := packet.Reader(msg)
			go protos.HandleRequest(hostid, reader, output)
		case msg := <-forward:
			SendChan(0, msg, output)
		}
	}

}

//---------------------------------------------------------- writer routine
func _write_routine(output chan []byte, conn net.Conn) {
	for {
		msg, ok := <-output
		if !ok {
			break
		}

		_, err := conn.Write(msg) // write operation is assumed to be atomic
		if err != nil {
			log.Println("Error send reply to GS:", err)
		}
	}
}
