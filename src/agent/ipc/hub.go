package ipc

import (
	"errors"
	"log"
	"net"
	"os"
)

import (
	"cfg"
	"hub/protos"
	"misc/packet"
	"sync"
)

var _conn net.Conn

func DialHub() {
	config := cfg.Get()

	conn, err := net.Dial("tcp", config["hub_service"])
	if err != nil {
		log.Println("Cannot connect to Hub")
		os.Exit(1)
	}

	_conn = conn
}

//--------------------------------------------------------- Send to Hub
func SendHub(id int32, tos int16, data []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.SendHub() failed")
		}
	}()

	// HUB protocol forwarding
	msg := protos.MSG{}
	msg.F_id = id
	msg.F_data = data
	return _send(packet.Pack(protos.Code["forward"], msg, nil), _conn)
}

var _seq_lock sync.Mutex

//---------------------------------------------------------- IPC send should be seqential.
func _send(data []byte, conn net.Conn) (err error) {
	_seq_lock.Lock()
	defer _seq_lock.Unlock()

	headwriter := packet.Writer()
	headwriter.WriteU16(uint16(len(data)))

	_, err = conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return
	}

	return nil
}
