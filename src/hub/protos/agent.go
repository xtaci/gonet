package protos

import (
	"log"
	"net"
	"sync/atomic"
	"sync"
)

import (
	"misc/packet"
)

const (
	MAXCHAN	= 65536
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
func HubAgent(in chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	MQ := make(chan []byte, MAXCHAN)

	_server_lock.Lock()
	_servers[hostid] = MQ
	_server_lock.Unlock()

	log.Printf("server id:%v connected\n", hostid)

	defer func() {
		_server_lock.Lock()
		delete(_servers, hostid)
		_server_lock.Unlock()

		log.Printf("server id:%v disconnected\n", hostid)
	}()

	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}

			if result := HandleRequest(hostid, msg); result != nil {
				_send(result, conn)
			}
		case msg, ok := <-MQ:
			if !ok {
				return
			}

			_send(msg, conn)
		}
	}

}

func _send(data []byte, conn net.Conn) {
	headwriter := packet.Writer()
	headwriter.WriteU16(uint16(len(data)))

	_, err := conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return
	}
}
