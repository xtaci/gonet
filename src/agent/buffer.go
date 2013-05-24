package agent

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

import (
	"cfg"
	"misc/packet"
)

var _BUFSIZE int32
var _MAXCHAN int32

type _RawPacket struct {
	data []byte // payload
}

type Buffer struct {
	ctrl    chan bool // receive exit signal
	pending chan *_RawPacket // pending Packet
	size    int32            // packet payload bytes count

	conn net.Conn // connection
}

//------------------------------------------------ Send packet
func (buf *Buffer) Send(data []byte) (err error) {
	if buf.size <= _BUFSIZE {
		rp := _RawPacket{data}
		buf.pending <- &rp

		atomic.AddInt32(&buf.size, int32(len(data)))
		return nil
	}

	return errors.New(fmt.Sprintf("Send Buffer Overflow, possible DoS attack. Remote: %v", buf.conn.RemoteAddr()))
}

//------------------------------------------------ packet sender goroutine
func (buf *Buffer) Start() {
	defer recover()

	for {
		select {
		case pkt := <-buf.pending:
			buf.raw_send(pkt)
			atomic.AddInt32(&buf.size, -int32(len(pkt.data)))
		case _ = <-buf.ctrl:
			return
		}
	}
}

//------------------------------------------------ send packet online
func (buf *Buffer) raw_send(pkt *_RawPacket) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(pkt.data)))
	writer.WriteRawBytes(pkt.data)

	_, err := buf.conn.Write(writer.Data())
	if err != nil {
		log.Println("Error send reply :", err)
		return
	}
}

//------------------------------------------------ create a new write buffer
func NewBuffer(conn net.Conn, ctrl chan bool) *Buffer {
	buf := Buffer{conn: conn, size: 0}
	buf.pending = make(chan *_RawPacket, _MAXCHAN)
	buf.ctrl = ctrl
	return &buf
}

func init() {
	_BUFSIZE = 65535

	config := cfg.Get()
	if config["write_buffer"] != "" {
		v, _ := strconv.Atoi(config["write_buffer"])
		_BUFSIZE = int32(v)
	}

	_MAXCHAN = _BUFSIZE / 16
}
