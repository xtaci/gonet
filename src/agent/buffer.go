package agent

import "net"
import "sync"
//import "time"
import "packet"
import "log"

type _RawPacket struct {
	Size uint16				// payload size
	SeqId uint32			// packet SEQ Num
	Data []byte				// payload
}

type Buffer struct {
	Pending chan _RawPacket	// Pending Packet
	Sent chan _RawPacket	// Sent but not ACKed 
	conn net.Conn			// connection
	_lock sync.Mutex		// for changing conn
	seqMax uint32			// Largest SEQ number
	ackMax uint32			// Largest ACK number
}

//---------------------------------------------------------Pend data
func (buf *Buffer) Send(data []byte) {
	buf.seqMax++
	rp := _RawPacket{Size:uint16(len(data)), SeqId:buf.seqMax, Data:data}
	buf.Pending <- rp
}

//---------------------------------------------------------Ack data
func (buf *Buffer) Ack(ackId uint32) {
	// retransit needed
	if ackId == buf.ackMax {
		buf.retransit()
	} else {
		//
		for {
			pkt := <-buf.Sent
			if buf.ackMax >= pkt.SeqId {
				break
			}
		}

		buf.ackMax = ackId
	}
}

//---------------------------------------------------------Change network-connection
func (buf *Buffer) ChangeConn(conn net.Conn) {
	buf._lock.Lock()
	defer buf._lock.Unlock()
	buf.conn = conn
}

//---------------------------------------------------------packet sender goroutine
func (buf *Buffer) Start() {
	for {
		pkt := <-buf.Pending
		buf.raw_send(pkt)
		buf.Sent <- pkt
	}
}

//---------------------------------------------------------retransit
func (buf *Buffer) retransit() {
	Sent2 := make(chan _RawPacket, 128)

	// retrieve all 'sent' packet & send again
	for {
		select {
		case pkt := <-buf.Sent:
			buf.raw_send(pkt)
			Sent2 <- pkt
		default:
			break;
		}
	}

	// push back into 'sent' channel
	for {
		select {
		case pkt := <-Sent2:
			buf.Sent <- pkt
		default:
			break;
		}
	}
}

func (buf *Buffer) raw_send(pkt _RawPacket) error {
	headwriter:= packet.PacketWriter()
	headwriter.WriteU16(uint16(len(pkt.Data)+8))
	headwriter.WriteU32(pkt.SeqId)
	headwriter.WriteU32(buf.ackMax)

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
	buf := Buffer{conn:conn, seqMax:0}
	buf.Pending = make(chan _RawPacket, 128)
	buf.Sent = make(chan _RawPacket, 128)
	return &buf
}
