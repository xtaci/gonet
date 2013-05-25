package ipc

import (
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
	"cfg"
	"misc/packet"
	. "types"
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

		for k, v := range seq_id {
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

		if seqval == 0 { // packet forwarding, deliver to MQ
			reader := packet.Reader(data)
			dest_id, err := reader.ReadS32()
			if err != nil {
				log.Println("forward: read dest_id failed.")
				goto L
			}

			sess := QueryOnline(dest_id)
			if sess == nil {
				log.Println("forward: user is offline.")
			} else {
				func() {
					defer func() {
						if x := recover(); x != nil {
							log.Println("forward: deliver to MQ failed.")
						}
					}()
					decoded := &IPCObject{}
					err := json.Unmarshal(data[reader.Pos():], decoded)
					if err != nil {
						log.Println("unable to decode forwared-IPC request")
					} else {
						sess.MQ <- *decoded
					}
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
		log.Println("Error send packet to HUB:", err)
		return nil
	}

	// wait ack
	ACK := make(chan []byte)
	_wait_ack_lock.Lock()
	_wait_ack[seq_id] = ACK
	_wait_ack_lock.Unlock()

	select {
	case msg := <-ACK:
		return msg[2:] // ignore protocol header
	case _ = <-time.After(10 * time.Second):
	}

	return nil
}

func init() {
	_wait_ack = make(map[uint64]chan []byte)
}
