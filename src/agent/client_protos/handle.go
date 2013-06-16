package protos

import (
	"misc/packet"
	. "types"
)

func P_heart_beat_req(sess *Session, reader *packet.Packet) []byte {
	// nothing should be done
	return nil
}

func P_user_login_req(sess *Session, reader *packet.Packet) []byte {
	//tbl := PKT_user_login_info(reader)
	ret := command_result_pack{}
	return packet.Pack(Code["user_login_succeed_ack"], ret, nil)
}
