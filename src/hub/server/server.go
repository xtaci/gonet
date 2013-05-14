package protos

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

import (
	"misc/packet"
)

const (
	MAXCHAN = 65536
)

//----------------------------------------------- logical game server chans
var _host_genid int32
var _servers map[int32]chan []byte
var _server_lock sync.RWMutex

func init() {
	log.SetPrefix("[HUB]")
	_servers = make(map[int32]chan []byte)
}

//------------------------------------------------ Hub processing
func HubAgent(incoming chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	// forward buffer
	forward := make(chan[]byte, MAXCHAN)
	// output buffer
	output := make(chan[]byte, MAXCHAN)

	_server_lock.Lock()
	_servers[hostid] = forward // message chan for forwarding to client
	_server_lock.Unlock()

	log.Printf("server id:%v connected\n", hostid)

	go _write_routine(output, conn)

	defer func() {
		_server_lock.Lock()
		delete(_servers, hostid)
		_server_lock.Unlock()

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
			go HandleRequest(hostid,reader,output)
		case msg := <-forward:
			go _send(0, msg, output)
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

//--------------------------------------------------------- send to Game Server
func _send(seqid uint64, data []byte, output chan[]byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data))+8)
	writer.WriteU64(seqid)		// piggyback seq id
	writer.WriteRawBytes(data)
	output <- writer.Data()
}
