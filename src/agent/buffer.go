package agent

import "net"
import "sync"
import "sync/atomic"

//import "time"
import "packet"
import "log"
import "errors"

const BUFSIZE = 65535
const MAXCHAN = 4096

type _RawPacket struct {
	Size  uint16 // payload size
	SeqId uint32 // packet SEQ Num
	Data  []byte // payload
}

type Buffer struct {
	Pending chan *_RawPacket // Pending Packet
	Sent    chan *_RawPacket // Sent but not ACKed 
	seqMax  uint32           // Largest SEQ number
	ackMax  uint32           // Largest ACK number

	conn      net.Conn   // connection
	byteCount int32      // packet payload bytes count
	_lock     sync.Mutex // for changing conn
}

//--------------------------------------------------------- Pend data
func (buf *Buffer) Send(data []byte) error {

	if buf.byteCount < BUFSIZE {
		buf.seqMax++
		rp := _RawPacket{Size: uint16(len(data)), SeqId: buf.seqMax, Data: data}
		buf.Pending <- &rp

		atomic.AddInt32(&buf.byteCount, int32(len(data)))
		return nil
	}

	return errors.New("Send Buffer Overflow, send rejected, possible DoS attack.")
}

//--------------------------------------------------------- Ack data
func (buf *Buffer) Ack(seqId uint32) {
	// retransit needed
	if seqId == buf.ackMax {
		buf.retransit()
	} else {
		// drop Acked packet
		for {
			pkt_old := <-buf.Sent
			if buf.ackMax >= pkt_old.SeqId {
				break
			}
		}

		// update max ack number
		atomic.StoreUint32(&buf.ackMax, seqId)
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
		pkt := <-buf.Pending
		buf.raw_send(pkt)
		buf.Sent <- pkt
		atomic.AddInt32(&buf.byteCount, -int32(len(pkt.Data)))
	}
}

//--------------------------------------------------------- retransit
func (buf *Buffer) retransit() {
	Sent2 := make(chan *_RawPacket, MAXCHAN)

	// retrieve all 'sent' packet & send again
	for {
		select {
		case pkt := <-buf.Sent:
			buf.raw_send(pkt)
			Sent2 <- pkt
		default:
			break
		}
	}

	// push back into 'sent' channel
	for {
		select {
		case pkt := <-Sent2:
			buf.Sent <- pkt
		default:
			break
		}
	}
}

func (buf *Buffer) raw_send(pkt *_RawPacket) error {
	headwriter := packet.PacketWriter()
	headwriter.WriteU16(uint16(len(pkt.Data) + 8))
	headwriter.WriteU32(pkt.SeqId)
	headwriter.WriteU32(atomic.LoadUint32(&buf.ackMax))

	buf._lock.Lock()
	defer buf._lock.Unlock()
	_, err := buf.conn.Write(headwriter.Data())
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return err
	}

	_, err = buf.conn.Write(pkt.Data)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}

func NewBuffer(conn net.Conn) *Buffer {
	buf := Buffer{conn: conn, seqMax: 0, byteCount: 0}
	buf.Pending = make(chan *_RawPacket, MAXCHAN)
	buf.Sent = make(chan *_RawPacket, MAXCHAN)
	return &buf
}
