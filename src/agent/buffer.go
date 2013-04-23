package agent

import (
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

import (
	"misc/packet"
)

const BUFSIZE = 65535
const MAXCHAN = 4096

type _RawPacket struct {
	size   uint16 // payload size
	seq_id uint32 // packet SEQ Num
	data   []byte // payload
}

type Buffer struct {
	ch_pending chan *_RawPacket // ch_pending Packet
	ch_sent    chan *_RawPacket // ch_sent but not ACKed 
	seq_max    uint32           // Largest SEQ number
	ack_max    uint32           // Largest ACK number

	conn     net.Conn   // connection
	byte_cnt int32      // packet payload bytes count
	_lock    sync.Mutex // for changing conn
}

//--------------------------------------------------------- Pend data
func (buf *Buffer) Send(data []byte) error {

	if buf.byte_cnt < BUFSIZE {
		buf.seq_max++
		rp := _RawPacket{size: uint16(len(data)), seq_id: buf.seq_max, data: data}
		buf.ch_pending <- &rp

		atomic.AddInt32(&buf.byte_cnt, int32(len(data)))
		return nil
	}

	return errors.New("Send Buffer Overflow, send rejected, possible DoS attack.")
}

//--------------------------------------------------------- Ack data
func (buf *Buffer) Ack(seqId uint32) {
	// retransit needed
	if seqId == buf.ack_max {
		buf.retransit()
	} else {
		// drop Acked packet
		for {
			pkt_old := <-buf.ch_sent
			if buf.ack_max >= pkt_old.seq_id {
				break
			}
		}

		// update max ack number
		atomic.StoreUint32(&buf.ack_max, seqId)
	}
}

//--------------------------------------------------------- Change network-connection
func (buf *Buffer) ChangeConn(conn net.Conn) {
	buf._lock.Lock()
	defer buf._lock.Unlock()
	buf.conn = conn
}

//--------------------------------------------------------- packet sender goroutine
func (buf *Buffer) Start() {
	for {
		pkt := <-buf.ch_pending
		buf.raw_send(pkt)
		buf.ch_sent <- pkt
		atomic.AddInt32(&buf.byte_cnt, -int32(len(pkt.data)))
	}
}

//--------------------------------------------------------- retransit
func (buf *Buffer) retransit() {
	ch_sent2 := make(chan *_RawPacket, MAXCHAN)

	// retrieve all 'sent' packet & send again
	for {
		select {
		case pkt := <-buf.ch_sent:
			buf.raw_send(pkt)
			ch_sent2 <- pkt
		default:
			break
		}
	}

	// push back into 'sent' channel
	for {
		select {
		case pkt := <-ch_sent2:
			buf.ch_sent <- pkt
		default:
			break
		}
	}
}

func (buf *Buffer) raw_send(pkt *_RawPacket) error {
	headwriter := packet.PacketWriter()
	headwriter.WriteU16(uint16(len(pkt.data) + 8))
	headwriter.WriteU32(pkt.seq_id)
	headwriter.WriteU32(atomic.LoadUint32(&buf.ack_max))

	buf._lock.Lock()
	defer buf._lock.Unlock()
	_, err := buf.conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return err
	}

	_, err = buf.conn.Write(pkt.data)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}

func NewBuffer(conn net.Conn) *Buffer {
	buf := Buffer{conn: conn, seq_max: 0, byte_cnt: 0}
	buf.ch_pending = make(chan *_RawPacket, MAXCHAN)
	buf.ch_sent = make(chan *_RawPacket, MAXCHAN)
	return &buf
}
