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
	list := ranklist.GetRankList(1, ranklist.Count())
	out := rank_list{}
	out.F_items = make([]rank_list_item, len(list))

	for i := range list {
		user := list[i]
		out.F_items[i].F_id = user.Id
		out.F_items[i].F_name = user.Name
		out.F_items[i].F_rank = user.Score
		out.F_items[i].F_state = int32(user.State)

		t := int32(user.ProtectTime.Unix() - time.Now().Unix())
		if t > 0 {
			out.F_items[i].F_protect_time = t
		} else {
			out.F_items[i].F_protect_time = 0
		}
	}

	writer := packet.Writer()
	return pack(Code["rank_list_ack"], out, writer), nil
}
