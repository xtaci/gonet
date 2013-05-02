package protos

import (
	. "types"
	"misc/packet"
	"hub/ranklist"
)

func _atk_player_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_command_id_pack(reader)
	writer := packet.Writer()
	success := user_snapshot{}
	failed := command_result_pack{}

	opponent := ranklist.Find(tbl.F_id)

	switch opponent.State {
	case ONLINE:
		failed.F_rst = 1
		ret = pack(Code["atk_player_faild_ack"], failed, writer)
	case BEING_RAID:
		failed.F_rst = 2
		ret = pack(Code["atk_player_faild_ack"], failed, writer)
	case PROTECTED:
		failed.F_rst = 3
		ret = pack(Code["atk_player_faild_ack"], failed, writer)
	default:
	}

	// 
	_fill_user_snapshot(&sess.User,&success)
	ret = pack(Code["atk_player_succeed_ack"], success, writer)
	return
}
