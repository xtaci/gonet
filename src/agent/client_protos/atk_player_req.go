package protos

import (
	"db/user_tbl"
	"agent/ipc"
	"misc/packet"
	. "types"
)

func P_atk_player_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := PKT_command_id_pack(reader)
	writer := packet.Writer()
	success := user_snapshot{}
	failed := command_result_pack{}

	if ipc.Raid(tbl.F_id) {
		opponent, e := user_tbl.Load(tbl.F_id)
		if e == nil {
			_fill_user_snapshot(&opponent, &success)
			ret = packet.Pack(Code["atk_player_succeed_ack"], success, writer)
			return
		}
	}

	// 
	info, err := ipc.GetInfo(tbl.F_id)
	failed.F_rst = info.Id
	ret = packet.Pack(Code["atk_player_faild_ack"], failed, writer)
	return
}
