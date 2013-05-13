package protos

import (
	"misc/packet"
	. "types"
)

func P_atk_monster_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := PKT_command_id_pack(reader)
	writer := packet.Writer()
	payload := command_result_pack{}

	//
	println(tbl.F_id)
	ret = packet.Pack(Code["atk_monster_ack"], payload, writer)
	return
}
