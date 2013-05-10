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

func init() {
	log.SetPrefix("[HUB]")
}

//--------------------------------------------------------- Hub processing
func HubAgent(in chan []byte, conn net.Conn) {
	hostid := atomic.AddInt32(&_host_genid, 1)
	MQ := make(chan []byte, MAXCHAN)

	_lock.Lock()
	_servers[hostid] = MQ
	_lock.Unlock()

	log.Println("server id:%v connected", hostid)

	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}

			if result := HandleRequest(hostid, msg); result != nil {
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
		//case msg, ok := <-MQ:
		}
	}

	_lock.Lock()
	delete(_servers, hostid)
	_lock.Unlock()
}


