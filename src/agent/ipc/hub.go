package ipc

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

import (
	"cfg"
	"hub/protos"
	"misc/packet"
	"sync"
	"sync/atomic"
)

var _conn net.Conn

//----------------------------------------------- connect to hub
func DialHub() {
	log.Println("Connecting to HUB")
	config := cfg.Get()

	conn, err := net.Dial("tcp", config["hub_service"])
	if err != nil {
		log.Println("Cannot connect to Hub")
		os.Exit(1)
	}

	_conn = conn

	log.Println("HUB connected")
	go HubReceiver(conn)
}

//----------------------------------------------- receive message from hub
func HubReceiver(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	seq_id := make([]byte, 4)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receving header:", err)
			break
		}

		// packet seq_id uint32
		n, err = io.ReadFull(conn, seq_id)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receving seq_id:", err)
			break
		}

		seqval := uint32(seq_id[0])<<24 | uint32(seq_id[1])<<16 | uint32(seq_id[2])<<8 | uint32(seq_id[3])

		// data
		size := int(header[0])<<8 | int(header[1])
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receving msg:", err)
			break
		}

		_wait_ack_lock.Lock()
		if ack, ok := _wait_ack[seqval]; ok {
			ack <- data
			delete(_wait_ack, seqval)
		} else {
			log.Println("Illegal packet sequence number from HuB")
		}
		_wait_ack_lock.Unlock()

	}
}

//------------------------------------------------ Forward to Hub
func ForwardHub(id int32, data []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New(x.(string))
		}
	}()

	// HUB protocol forwarding
	msg := protos.MSG{}
	msg.F_id = id
	msg.F_data = data
	ack := _call(packet.Pack(protos.Code["forward"], msg, nil), _conn)
	if ack != nil {
		panic("ForwardHub failed or timed-out")
	}

	return nil
}

// send lock
var _seq_lock sync.Mutex

// packet sequence number generator
var _seq_id uint32

// waiting ACK queue.
var _wait_ack map[uint32]chan []byte
var _wait_ack_lock sync.Mutex

//------------------------------------------------ IPC send should be seqential
func _call(data []byte, conn net.Conn) (ret []byte) {
	seq_id := atomic.AddUint32(&_seq_id, 1)

	_seq_lock.Lock()
	headwriter := packet.Writer()
	headwriter.WriteU16(uint16(len(data)) + 4) // data + seq id
	headwriter.WriteU32(seq_id)

	_, err := conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send packet header:", err.Error())
		return nil
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error send packet data:", err.Error())
		return nil
	}
	_seq_lock.Unlock()

	// wait ack
	ACK := make(chan []byte)
	_wait_ack_lock.Lock()
	_wait_ack[seq_id] = ACK
	_wait_ack_lock.Unlock()

	select {
	case msg := <-ACK:
		return msg
	case _ = <-time.After(10 * time.Second):
	}

	return nil
}

func init() {
	_wait_ack = make(map[uint32]chan []byte)
}
