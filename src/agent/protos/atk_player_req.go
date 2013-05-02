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

	state := ranklist.GetState(tbl.F_id)
	switch state {
	case ONLINE, BEING_RAID, PROTECTED:
		failed.F_rst = state
		ret = pack(Code["atk_player_faild_ack"], failed, writer)
	default:
		if ranklist.ChangeState(tbl.F_id, int32(FREE), int32(BEING_RAID)) {
			opponent := ranklist.GetUserCopy(tbl.F_id)
			_fill_user_snapshot(&opponent, &success)
			ret = pack(Code["atk_player_succeed_ack"], success, writer)
		} else {
			failed.F_rst = ranklist.GetState(tbl.F_id)
			ret = pack(Code["atk_player_faild_ack"], failed, writer)
		}
	}

	// 
	return
}
