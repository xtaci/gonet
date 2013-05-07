package protos

import (
	"db/user_tbl"
	"hub/ranklist"
	"misc/packet"
	. "types"
)

func _atk_player_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_command_id_pack(reader)
	writer := packet.Writer()
	success := user_snapshot{}
	failed := command_result_pack{}

	state := ranklist.GetState(tbl.F_id)
	if state == ranklist.FREE {
		if ranklist.ChangeState(tbl.F_id, state, int32(ranklist.BEING_RAID)) {
			opponent, e := user_tbl.Read(tbl.F_id)
			if e == nil {
				_fill_user_snapshot(&opponent, &success)
				ret = pack(Code["atk_player_succeed_ack"], success, writer)
				return
			}
		}
	}

	// 
	failed.F_rst = ranklist.GetState(tbl.F_id)
	ret = pack(Code["atk_player_faild_ack"], failed, writer)
	return
}
