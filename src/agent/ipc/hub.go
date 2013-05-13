package ipc

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
	"sync"
	"sync/atomic"
)

import (
	"cfg"
	"hub/server"
	"misc/packet"
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
	seq_id := make([]byte, 8)

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

		seqval := uint64(0)

		for k,v := range seq_id {
			seqval |= uint64(v) << uint((7-k)*8)
		}

		// data
		size := int(header[0])<<8 | int(header[1]) - 8
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receving msg:", err)
			break
		}

		if seqval == 0 {		// packet forwarding, deliver to MQ
			reader := packet.Reader(data)
			forward_id, err := reader.ReadS32()
			if err != nil {
				log.Println("packet forwarding error")
				goto L
			}

			sess := Query(forward_id)
			if sess == nil {
				log.Println("forward failed, maybe user is offline?")
			} else {
				func() {
					defer func() {
						if x := recover(); x != nil {
							log.Println("forward to MQ failed, the user is so lucky")
						}
					}()
					sess.MQ <- data[reader.Pos():]	// the payload is the message
				}()
			}
		} else {
			_wait_ack_lock.Lock()
			if ack, ok := _wait_ack[seqval]; ok {
				ack <- data
				delete(_wait_ack, seqval)
			} else {
				log.Println("Illegal packet sequence number from HUB")
			}
			_wait_ack_lock.Unlock()
		}
L:
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
	ack := _call(packet.Pack(protos.Code["forward"], msg, nil))
	if ack != nil {
		panic("ForwardHub failed or timed-out")
	}

	return nil
}

// send lock
var _seq_lock sync.Mutex

// packet sequence number generator
var _seq_id uint64

// waiting ACK queue.
var _wait_ack map[uint64]chan []byte
var _wait_ack_lock sync.Mutex

//------------------------------------------------ IPC send should be seqential
func _call(data []byte) (ret []byte) {
	seq_id := atomic.AddUint64(&_seq_id, 1)

	_seq_lock.Lock()
	headwriter := packet.Writer()
	headwriter.WriteU16(uint16(len(data)) + 8) // data + seq id
	headwriter.WriteU64(seq_id)

	_, err := _conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send packet header:", err.Error())
		return nil
	}

	_, err = _conn.Write(data)
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
	_wait_ack = make(map[uint64]chan []byte)
}
