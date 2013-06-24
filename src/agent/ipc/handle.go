package ipc

import (
	"agent/AI"
	"db/user_tbl"
	"misc/packet"
	. "types"
)

func P_heart_beat_req(sess *Session, reader *packet.Packet) []byte {
	// nothing should be done
	return nil
}

func P_user_login_req(sess *Session, reader *packet.Packet) []byte {
	tbl, _ := PKT_user_login_info(reader)
	ret := command_result_pack{}

	if user := user_tbl.LoginMac(tbl.F_user_name, tbl.F_mac_addr); user != nil {
		sess.User = user
		AI.LoginProc(sess)
	}

	return packet.Pack(Code["user_login_succeed_ack"], &ret, nil)
}

func P_talk_req(sess *Session, reader *packet.Packet) []byte {
	tbl, _ := PKT_talk(reader)
	dest := user_tbl.Query(tbl.F_user)
	Send(sess.User.Id, dest.Id, SERVICE_TALK, false, tbl.F_msg)
	return nil
}
