package protos

import . "types"
import "misc/packet"

func _atk_player_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_command_id_pack(reader)
	writer := packet.Writer()
	success := user_snapshot{}
	//failed := command_result_pack{}

	// 
	println(tbl.F_id)
	_fill_user_snapshot(&sess.User,&success)
	ret = pack(Code["atk_player_succeed_ack"], success, writer)
	return
}
