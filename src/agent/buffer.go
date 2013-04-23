package agent

import "net"
import "sync"
//import "time"
import "packet"
import "log"

type _RawPacket struct {
	Size uint16
	SeqId uint32
	Data []byte
}

type Buffer struct {
	Pending chan _RawPacket	// CSP
	Sent chan _RawPacket
	num_packet uint32
	conn net.Conn
	_lock sync.Mutex
	seqMax uint32
	ackMax uint32
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

func (buf *Buffer) Start() {
	go buf.retransit()
	go buf.sender()
}

//---------------------------------------------------------retransit
func (buf *Buffer) retransit() {
	Sent2 := make(chan _RawPacket, 128)

	for {
		select {
		case pkt := <-buf.Sent:
			buf.send(pkt)
			Sent2 <- pkt
		default:
			break;
		}
	}

	buf.Sent = Sent2
}

//---------------------------------------------------------packet sender
func (buf *Buffer) sender() {
	for {
		pkt := <-buf.Pending
		buf.send(pkt)
		buf.Sent <- pkt
	}
}

func (buf *Buffer) send(pkt _RawPacket) error {
	writer := packet.PacketWriter()
	writer.WriteU16(uint16(len(pkt.Data)+8))
	writer.WriteU32(pkt.SeqId)
	writer.WriteU32(buf.ackMax)

	_, err := buf.conn.Write(writer.Data())
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
	buf := Buffer{conn:conn,num_packet:0, seqMax:0}
	buf.Pending = make(chan _RawPacket, 128)
	buf.Sent = make(chan _RawPacket, 128)
	return &buf
}
