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
	seqId uint32
	ackId uint32
}

//---------------------------------------------------------Pend data
func (buf *Buffer) Send(data []byte) {
	buf.seqId++
	rp := _RawPacket{Size:uint16(len(data)), SeqId:buf.seqId, Data:data}
	buf.Pending <- rp
}

//---------------------------------------------------------Ack data
func (buf *Buffer) Ack(ackId uint32) {
	// first ack
	for {
		pkt := <-buf.Sent
		if ackId >= pkt.SeqId {
			break
		}
	}

	buf.ackId = ackId
}

func (buf *Buffer) Start() {
	go buf.retransit()
	go buf.sender()
}

//---------------------------------------------------------retransit
func (buf *Buffer) retransit() {
}

//---------------------------------------------------------packet sender
func (buf *Buffer) sender() {
	for {
		pkt := <-buf.Pending
		buf.send(pkt)
	}
}

func (buf *Buffer) send(pkt _RawPacket) error {
	writer := packet.PacketWriter()
	writer.WriteU16(uint16(len(pkt.Data)+8))
	writer.WriteU32(pkt.SeqId)
	writer.WriteU32(buf.ackId)

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
	buf := Buffer{conn:conn,num_packet:0, seqId:0}
	buf.Pending = make(chan _RawPacket, 128)
	buf.Sent = make(chan _RawPacket, 128)
	return &buf
}
