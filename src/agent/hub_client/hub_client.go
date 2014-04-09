package hub_client

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

import (
	"agent/gsdb"
	"cfg"
	"db/forward_tbl"
	. "helper"
	"misc/packet"
	. "types"
)

type _request struct {
	data []byte
	ack  chan []byte
}

type _response struct {
	seq_id uint32
	data   []byte
}

var (
	_conn   net.Conn
	_caller chan _request
	_callee chan _response
)

func init() {
	_caller = make(chan _request, 20000)
	_callee = make(chan _response, 50000)
	go _caller_routine()
}

//----------------------------------------------- connect to hub
func DialHub() {
RETRY:
	INFO("Connecting to HUB")
	config := cfg.Get()

	addr, err := net.ResolveTCPAddr("tcp", config["hub_service"])
	if err != nil {
		ERR(err, addr)
		panic("cannot read hub_service from config")
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		ERR("connect to hub failed", err, "waiting 5 seconds")
		time.Sleep(5 * time.Second)
		goto RETRY
	}

	// set parameter
	conn.SetLinger(-1)
	_conn = conn
	go _receiver(_conn)

	INFO("HUB connected")
}

func _caller_routine() {
	// waiting ACK queue.
	wait_ack := make(map[uint32]chan []byte)
	var seq_id uint32
	for {
		select {
		case req := <-_caller:
			seq_id++
			if seq_id == 0 {
				seq_id++
			}

			wait_ack[seq_id] = req.ack
			writer := packet.Writer()
			writer.WriteU16(uint16(len(req.data)) + 4) // data + seq id
			writer.WriteU32(seq_id)
			writer.WriteRawBytes(req.data)
			// send the packet
			n, err := _conn.Write(writer.Data())
			if err != nil {
				ERR("Error send packet to HUB, bytes:", n, "reason:", err)
				// reconnect
				DialHub()
			}
		case resp := <-_callee: // callee
			ack, ok := wait_ack[resp.seq_id]
			if !ok {
				ERR("Illegal(or expired) packet seqid:", resp.seq_id, "from HUB", resp.data)
				continue
			}
			ack <- resp.data
			close(ack)
			delete(wait_ack, resp.seq_id)
		}
	}
}

//------------------------------------------------ a single IPC call
func _call(data []byte) (ret []byte) {
	// packet size test
	if len(data) > packet.PACKET_LIMIT {
		ERR("PACKET SIZE TOO LARGE... CHANGE YOUR DESIGN!", len(data))
		return nil
	}

	start := time.Now()
	defer func() {
		end := time.Now()
		nr := int16(data[0])<<8 | int16(data[1])
		DEBUG("HUB RPC #", nr, "call time:", end.Sub(start))
	}()

	// construct a request
	ack := make(chan []byte, 1)
	req := _request{data, ack}
	_caller <- req

	// wait for response
	select {
	case msg := <-ack:
		return msg
	case <-time.After(10 * time.Second):
		ERR("HUB is not responding...")
	}

	return nil
}

//---------------------------------------------------------- deliver an IPCObject to a user
func _deliver(obj *IPCObject) {
	sess := gsdb.QueryOnline(obj.DestID)
	if sess != nil {
		func() {
			defer func() {
				if x := recover(); x != nil {
					forward_tbl.Push(obj)
				}
			}()

			sess.MQ <- *obj
		}()
	} else {
		forward_tbl.Push(obj)
	}
}

//----------------------------------------------- receive message from hub
// receiver() will automatically exit if something wrong with the connection
func _receiver(conn net.Conn) {
	defer func() {
		recover()
	}()

	header := make([]byte, 2)
	seqval := make([]byte, 4)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			ERR("error receiving header:", n, err)
			break
		}

		// packet seq_id
		n, err = io.ReadFull(conn, seqval)
		if err != nil {
			ERR("error receiving seq_id:", n, err)
			break
		}

		// read big-endian header
		seq_id := binary.BigEndian.Uint32(seqval)
		size := binary.BigEndian.Uint16(header) - 4
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			ERR("error receiving msg:", err)
			break
		}

		// two kinds of IPC:
		// a). Hub Sends to GS, sequence number is not required (set to 0), just forwarding to session
		// b). Call, sequence number is needed, send will wake up blocking-chan.
		if seq_id == 0 {
			obj := &IPCObject{}
			err := json.Unmarshal(data, obj)
			if err != nil {
				ERR("unable to decode received IPCObject")
				continue
			}

			_deliver(obj)
		} else {
			resp := _response{seq_id, data}
			_callee <- resp
		}
	}
}

func checkErr(err error) {
	if err != nil {
		ERR(err)
		panic("error occured in protocol module")
	}
}
