package protos

import (
	"time"
)

import (
	"agent/ipc"
	"misc/packet"
	. "types"
)

func P_rank_list_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	ids, scores, err := ipc.GetList(1, -1)

	if err != nil {
			return nil, err
	}

	out := rank_list{}
	out.F_items = make([]rank_list_item, len(ids))

	for k, v := range ids {
		info, _ := ipc.GetInfo(v)
		out.F_items[k].F_id = v
		out.F_items[k].F_name = info.Name
		out.F_items[k].F_rank = scores[k]
		out.F_items[k].F_state = info.State

		t := int32(info.ProtectTime - time.Now().Unix())
		if t > 0 {
			out.F_items[k].F_protect_time = int32(t)
		} else {
			out.F_items[k].F_protect_time = 0
		}
	}

	writer := packet.Writer()
	return packet.Pack(Code["rank_list_ack"], out, writer), nil
}
