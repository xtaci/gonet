package main

import (
	"log"
	"time"
)

import (
	"agent/ipc_service"
	"agent/net"
	. "helper"
	"misc/packet"
	. "types"
)

//----------------------------------------------- client protocol handle proxy
func UserRequestProxy(sess *Session, p []byte) []byte {
	defer PrintPanicStack()
	// decrypt
	if sess.Flag&SESS_ENCRYPT != 0 {
		sess.Decoder.Codec(p)
	}

	// encapsulate into reader
	reader := packet.Reader(p)

	// client timestamp check
	// mainly for REPLAY-ATTACK
	client_elapsed, err := reader.ReadU32()
	if err != nil {
		ERR("read client timestamp failed.", err)
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	client_time := sess.ConnectTime.Unix() + int64(client_elapsed)/1000
	now := time.Now().Unix()
	if client_time > now+PACKET_ERROR || client_time < now-PACKET_EXPIRE {
		ERR("client timestamp is illegal.", client_elapsed, client_time, now)
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// read protocol number
	b, err := reader.ReadS16()
	if err != nil {
		ERR("read protocol number failed.")
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// handle validation
	handle := net.ProtoHandler[b]
	if handle == nil {
		ERR("service id", b, "not bind")
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// before HOOK
	if !_before_hook(sess, b) {
		ERR("before hook failed, code", b)
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}

	// handle packet
	start := time.Now()
	ret := handle(sess, reader)
	end := time.Now()

	uid := int32(-1)
	name := ""
	if sess.Flag&SESS_LOGGED_IN != 0 {
		uid = sess.User.Id
		name = sess.User.Name
	}

	log.Printf("\033[0;36m[REQ] %v\tbytes[in:%v out:%v seq:%v]\tusr:[%v %v]\ttime:%v\033[0m\n", net.RCode[b], len(p)-6, len(ret), sess.PacketCount, uid, name, end.Sub(start))
	// after HOOK
	_after_hook(sess, net.RCode[b])
	sess.MarkDirty()
	return ret
}

//----------------------------------------------- IPC proxy
func IPCRequestProxy(sess *Session, p *IPCObject) []byte {
	defer PrintPanicStack()

	// handle validation
	handle := ipc_service.IPCHandler[p.Service]
	if handle == nil {
		ERR("ipc service", p.Service, "not bind, internal service error.")
		return nil
	}

	// process IPCObject
	start := time.Now()
	ret := handle(sess, p)
	end := time.Now()
	log.Printf("\033[0;36m[IPC] %v\t%v\033[0m\n", p.Service, end.Sub(start))

	sess.MarkDirty()
	return ret
}

//---------------------------------------------------------- after hook
func _after_hook(sess *Session, rcode string) {
	if sess.Flag&SESS_LOGGED_IN == 0 {
		return
	}
}

//---------------------------------------------------------- before hook
func _before_hook(sess *Session, rcode int16) bool {
	return true
}
