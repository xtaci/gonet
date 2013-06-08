package protos

import (
	"time"
)

import (
	"agent/AI"
	"misc/packet"
	. "types"
)

var EPOCH = time.Unix(0, 0)

func P_user_login_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	AI.LoginProc(sess)
	return
}
