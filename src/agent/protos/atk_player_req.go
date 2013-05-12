package protos

import (
	"db/user_tbl"
	"hub/ranklist"
	"misc/packet"
	. "types"
)

func P_atk_player_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := PKT_command_id_pack(reader)
	writer := packet.Writer()
	success := user_snapshot{}
	failed := command_result_pack{}

	if ranklist.Raid(tbl.F_id) {
		opponent, e := user_tbl.Load(tbl.F_id)
		if e == nil {
			_fill_user_snapshot(&opponent, &success)
			ret = packet.Pack(Code["atk_player_succeed_ack"], success, writer)
			return
		}
	}

	// 
	failed.F_rst = int32(ranklist.State(tbl.F_id))
	ret = packet.Pack(Code["atk_player_faild_ack"], failed, writer)
	return
}
