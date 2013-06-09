package protos

import (
	"fmt"
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
	fmt.Println(ret)
	return packet.Pack(Code["user_login_succeed_ack"], ret, nil), nil
}
