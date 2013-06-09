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

func P_user_login_req(sess *Session, reader *packet.Packet) ([]byte, error) {
	ret := command_result_pack{}
	AI.LoginProc(sess)
	return packet.Pack(Code["user_login_succeed_ack"], ret, nil), nil
}
