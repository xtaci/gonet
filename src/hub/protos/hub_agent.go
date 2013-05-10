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

// global comm
var _host_genid int32
var _servers map[int32]chan []byte
var _lock sync.RWMutex

//--------------------------------------------------------- Hub processing
func HubAgent(in chan []byte, conn net.Conn) {
	id := atomic.AddInt32(&_host_genid, 1)
	_lock.Lock()
	_servers[id] = make(chan []byte, MAXCHAN)
	_lock.Unlock()

	log.Println("server id:%v connected", id)

	for {
		msg, ok := <-in

		if !ok {
			return
		}

		if result := HandleRequest(msg); result != nil {
			headwriter := packet.Writer()
			headwriter.WriteU16(uint16(len(result)))

			_, err := conn.Write(headwriter.Data())
			if err != nil {
				log.Println("Error send reply header:", err.Error())
				return
			}

			_, err = conn.Write(result)
			if err != nil {
				log.Println("Error send reply msg:", err.Error())
				return
			}
		}
	}

	_lock.Lock()
	delete(_servers, id)
	_lock.Unlock()
}

func init() {
	log.SetPrefix("[HUB]")
}
