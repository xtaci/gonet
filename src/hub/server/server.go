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
			seqid, err := reader.ReadU64()	// read seqid 
			if err != nil {
				log.Println("Read Sequence Id failed.", err)
				continue
			}

			if result := HandleRequest(hostid, reader); result != nil {
				_send(seqid, result, conn)
			}
		case msg := <-forward:
			_send(0, msg, conn)
		}
	}

}

//--------------------------------------------------------- send to Game Server
func _send(seqid uint64, data []byte, conn net.Conn) {
	headwriter := packet.Writer()
	headwriter.WriteU16(uint16(len(data))+8)
	headwriter.WriteU64(seqid)		// piggyback seq id

	_, err := conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send reply header:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error send reply msg:", err)
		return
	}
}
