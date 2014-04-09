package main

import (
	"net"
	"strconv"
)

import (
	"cfg"
	. "helper"
	"misc/packet"
	. "types"
)

type Buffer struct {
	ctrl    chan bool   // receive exit signal
	pending chan []byte // pending Packet
	max     int         // max queue size
	conn    net.Conn    // connection
	sess    *Session    // session
}

//------------------------------------------------ send packet
// !IMPORTANT! once closed, never Send!!!
func (buf *Buffer) Send(data []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			WARN("buffer.Send failed", x)
		}
	}()

	if buf.sess.Flag&SESS_ENCRYPT != 0 { // if encryption has setup
		buf.sess.Encoder.Codec(data)
	} else if buf.sess.Flag&SESS_KEYEXCG != 0 { // whether we just exchanged the key
		buf.sess.Flag &= ^SESS_KEYEXCG
		buf.sess.Flag |= SESS_ENCRYPT
	}
	buf.pending <- data
	return nil
}

//------------------------------------------------ packet sender goroutine
func (buf *Buffer) Start() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in buffer goroutine", x)
		}
	}()

	for {
		select {
		case data := <-buf.pending:
			buf.raw_send(data)
		case <-buf.ctrl: // session end, send final data
			close(buf.pending)
			for data := range buf.pending {
				buf.raw_send(data)
			}
			// close connection
			buf.conn.Close()
			return
		}
	}
}

//------------------------------------------------ packet online
func (buf *Buffer) raw_send(data []byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)))
	writer.WriteRawBytes(data)

	//nr := int16(data[0])<<8 | int16(data[1])
	//log.Printf("\033[37;44m[ACK] %v\t%v\tSIZE:%v\033[0m\n", nr, proto.RCode[nr], len(data))
	n, err := buf.conn.Write(writer.Data())
	if err != nil {
		ERR("Error send reply, bytes:", n, "reason:", err)
		return
	}
}

//------------------------------------------------ create a new write buffer
func NewBuffer(sess *Session, conn net.Conn, ctrl chan bool) *Buffer {
	config := cfg.Get()
	max, err := strconv.Atoi(config["outqueue_size"])
	if err != nil {
		max = DEFAULT_OUTQUEUE_SIZE
		WARN("cannot parse outqueue_size from config", err, "using default:", max)
	}

	buf := Buffer{conn: conn}
	buf.sess = sess
	buf.pending = make(chan []byte, max)
	buf.ctrl = ctrl
	buf.max = max
	return &buf
}
