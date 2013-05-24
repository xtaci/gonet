package event_client

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"cfg"
	"misc/packet"
)

var _conn net.Conn

//----------------------------------------------- connect to Event server
func DialEvent() {
	log.Println("Connecting to Event server")
	config := cfg.Get()

	conn, err := net.Dial("tcp", config["event_service"])
	if err != nil {
		log.Println("Cannot connect to Event server")
		os.Exit(1)
	}

	_conn = conn

	log.Println("Event Service Connected")
	go EventReceiver(conn)
}

//----------------------------------------------- receive ack from Event Server
func EventReceiver(conn net.Conn) {
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
		for k, v := range seq_id {
			seqval |= uint64(v) << uint((7-k)*8)
		}

		// data
		size := int(header[0])<<8 | int(header[1]) - len(seq_id)
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receving msg:", err)
			break
		}

		// acks
		_wait_ack_lock.Lock()
		if ack, ok := _wait_ack[seqval]; ok {
			ack <- data
			delete(_wait_ack, seqval)
		} else {
			log.Println("Illegal packet sequence number from Event Server")
		}
		_wait_ack_lock.Unlock()
	}
}

// packet sequence number generator
var _seq_id uint64

// waiting ACK queue.
var _wait_ack map[uint64]chan []byte
var _wait_ack_lock sync.Mutex

//------------------------------------------------ IPC send should be seqential
func _call(data []byte) (ret []byte) {
	seq_id := atomic.AddUint64(&_seq_id, 1)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8) // data + seq id
	writer.WriteU64(seq_id)
	writer.WriteRawBytes(data)

	_, err := _conn.Write(writer.Data())
	if err != nil {
		log.Println("Error send packet to Event Server:", err)
		return nil
	}

	// wait ack
	ACK := make(chan []byte)
	_wait_ack_lock.Lock()
	_wait_ack[seq_id] = ACK
	_wait_ack_lock.Unlock()

	select {
	case msg := <-ACK:
		return msg[2:]	// the first 2-bytes are TYPE, just ignore it
	case _ = <-time.After(10 * time.Second):
	}

	return nil
}

func init() {
	_wait_ack = make(map[uint64]chan []byte)
}
