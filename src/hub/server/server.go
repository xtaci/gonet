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
	forward := make(chan[]byte)

	_server_lock.Lock()
	_servers[hostid] = forward // message chan for forwarding to client
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
		case msg, ok := <-incoming:
			if !ok {
				return
			}

			reader := packet.Reader(msg)
			go HandleRequest(hostid,reader,conn)
		case msg := <-forward:
			go _send(0, msg, conn)
		}
	}

}

