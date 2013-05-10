package protos

import (
	"time"
)

import (
	"hub/ranklist"
	"misc/packet"
	. "types"
)

func _rank_list_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	id, score := ranklist.GetList(1, -1)
	out := rank_list{}
	out.F_items = make([]rank_list_item, len(id))

	for k,v := range id {
		out.F_items[k].F_id = v
		out.F_items[k].F_name = ranklist.Name(v)
		out.F_items[k].F_rank = score[k]
		out.F_items[k].F_state = int32(ranklist.State(v))

		t := int32(ranklist.ProtectTime(v) - time.Now().Unix())
		if t > 0 {
			out.F_items[k].F_protect_time = t
		} else {
			out.F_items[k].F_protect_time = 0
		}
	}

	writer := packet.Writer()
	return packet.Pack(Code["rank_list_ack"], out, writer), nil
}
