package protos

import . "types"
import "misc/packet"

func _atk_monster_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_command_id_pack(reader)
	writer := packet.Writer()
	payload := command_result_pack{}

	//
	println(tbl.F_id)
	ret = pack(Code["atk_monster_ack"], payload, writer)
	return
}
