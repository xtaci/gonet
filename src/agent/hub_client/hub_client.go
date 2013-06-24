package hub_client

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"agent/gsdb"
	"cfg"
	"db/forward_tbl"
	"helper"
	"misc/packet"
	. "types"
)

var _conn net.Conn

//----------------------------------------------- connect to hub
func DialHub() {
	log.Println("Connecting to HUB")
	config := cfg.Get()

	addr, err := net.ResolveTCPAddr("tcp", config["hub_service"])
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

		// two kinds of IPC:
		// a). Hub Sends to GS, sequence number is not required (set to 0), just forwarding to session
		// b). Call, sequence number is needed, send will wake up blocking-chan.
		//
		if seqval == 0 {
			obj := &IPCObject{}
			err := json.Unmarshal(data, obj)
			if err != nil {
				log.Println("unable to decode received IPCObject")
				continue
			}

			sess := gsdb.QueryOnline(obj.DestID)
			if sess == nil {
				// if the user is disconnected
				forward_tbl.Push(obj)
			} else {
				func() {
					defer func() {
						if x := recover(); x != nil {
							log.Println("forward: deliver to MQ failed.")
							forward_tbl.Push(obj)
						}
					}()

					sess.MQ <- *obj
				}()
			}
		} else {
			_wait_ack_lock.Lock()
			if ack, ok := _wait_ack[seqval]; ok {
				ack <- data
				delete(_wait_ack, seqval)
			} else {
				log.Printf("Illegal packet sequence number [%x] from HUB", seqval)
			}
			_wait_ack_lock.Unlock()
		}
	}
}

// packet sequence number generator
// assume we have 100,000 messages per second, _seq_id will overflow in
// 2^64 / 100000 / 3600 / 24 / 365 = 5.8 million years
var _seq_id uint64

// waiting ACK queue.
var _wait_ack map[uint64]chan []byte
var _wait_ack_lock sync.Mutex

//------------------------------------------------ IPC call
func _call(data []byte) (ret []byte) {
	// packet creation
	seq_id := atomic.AddUint64(&_seq_id, 1)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8) // data + seq id
	writer.WriteU64(seq_id)
	writer.WriteRawBytes(data)

	// add seq_id to waiting queue
	ACK := make(chan []byte)
	_wait_ack_lock.Lock()
	_wait_ack[seq_id] = ACK
	_wait_ack_lock.Unlock()

	// send the packet
	_, err := _conn.Write(writer.Data())
	if err != nil {
		log.Println("Error send packet to HUB:", err)
		_wait_ack_lock.Lock()
		delete(_wait_ack, seq_id)
		_wait_ack_lock.Unlock()
		return nil
	}

	select {
	case msg := <-ACK:
		return msg
	case <-time.After(10 * time.Second):
		log.Println("HUB is not responding...")
	}

	return nil
}

func init() {
	_wait_ack = make(map[uint64]chan []byte)
}
