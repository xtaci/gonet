package stats_client

import (
	"encoding/binary"
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
	"helper"
	"misc/packet"
)

var _conn net.Conn

//----------------------------------------------- connect to Stats server
func DialStats() {
	log.Println("Connecting to Stats server")
	config := cfg.Get()

	addr, err := net.ResolveTCPAddr("tcp", config["stats_service"])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	helper.SetConnParam(conn)
	_conn = conn

	log.Println("Stats Service Connected")
	go StatsReceiver(conn)
}

//----------------------------------------------- receive ack from Stats Server
func StatsReceiver(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	seq_id := make([]byte, 8)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receiving header:", err)
			break
		}

		// packet seq_id uint32
		n, err = io.ReadFull(conn, seq_id)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receiving seq_id:", err)
			break
		}

		// read big-endian header
		seqval := binary.BigEndian.Uint64(seq_id)
		size := binary.BigEndian.Uint16(header) - 8
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receiving msg:", err)
			break
		}

		// acks
		_wait_ack_lock.Lock()
		if ack, ok := _wait_ack[seqval]; ok {
			ack <- data
			delete(_wait_ack, seqval)
		} else {
			log.Println("Illegal packet sequence number from Stats Server")
		}
		_wait_ack_lock.Unlock()
	}
}

// packet sequence number generator
var _seq_id uint64

// waiting ACK queue.
var _wait_ack map[uint64]chan []byte
var _wait_ack_lock sync.Mutex

//------------------------------------------------ call remote function
func _call(data []byte) (ret []byte) {
	seq_id := atomic.AddUint64(&_seq_id, 1)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8) // data + seq id
	writer.WriteU64(seq_id)
	writer.WriteRawBytes(data)

	// wait ack
	ACK := make(chan []byte)
	_wait_ack_lock.Lock()
	_wait_ack[seq_id] = ACK
	_wait_ack_lock.Unlock()

	// send the packet
	_, err := _conn.Write(writer.Data())
	if err != nil {
		log.Println("Error send packet to Stats Server:", err)
		_wait_ack_lock.Lock()
		delete(_wait_ack, seq_id)
		_wait_ack_lock.Unlock()
		return nil
	}

	select {
	case msg := <-ACK:
		return msg
	case <-time.After(10 * time.Second):
		log.Println("StatsServer is not responding...")
	}

	return nil
}

func init() {
	_wait_ack = make(map[uint64]chan []byte)
}
