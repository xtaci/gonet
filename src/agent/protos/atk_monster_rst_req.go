package protos

import . "types"
import "misc/packet"

func _atk_monster_rst_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_atk_monster_rst_req(reader)
	println(tbl.F_protect_time)
	return
}
