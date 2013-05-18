package main

import (
	"log"
	"net"
	"sync/atomic"
)

import (
	"misc/packet"
	"hub/protos"
)

const (
	MAXCHAN = 65536
)

//----------------------------------------------- logical game server chans
var _host_genid int32

func init() {
	log.SetPrefix("[HUB]")
}

//--------------------------------------------------------- send
func _send(seqid uint64, data []byte, output chan[]byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data))+8)
	writer.WriteU64(seqid)		// piggyback seq id
	writer.WriteRawBytes(data)
	output <- writer.Data()
}

//------------------------------------------------ Hub processing
func HubAgent(incoming chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	// forward buffer
	forward := make(chan[]byte, MAXCHAN)
	// output buffer
	output := make(chan[]byte, MAXCHAN)

	protos.ServerLock.Lock()
	protos.Servers[hostid] = forward // message chan for forwarding to client
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

	for {
		select {
		case msg, ok := <-incoming:
			if !ok {
				return
			}

			reader := packet.Reader(msg)
			go protos.HandleRequest(hostid,reader,output)
		case msg := <-forward:
			_send(0, msg, output)
		}
	}

}

//----------------------------------------------- to gs write buffer 
func _write_routine(output chan[]byte, conn net.Conn) {
	for {
		msg, ok := <-output
		if !ok {
			break
		}

		_, err := conn.Write(msg)	// write operation is assumed to be atomic
		if err != nil {
			log.Println("Error send reply to GS:", err)
		}
	}
}


