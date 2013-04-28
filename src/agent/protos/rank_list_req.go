package protos

import (
	"time"
)

import (
	. "types"
	"misc/packet"
	"agent/ipc"
	"hub/ranklist"
	"fmt"
	"db/user_tbl"
)

func _rank_list_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	list := ranklist.GetRankList(1, ranklist.Count())
	out := rank_list{}
	out.F_items = make([]rank_list_item, len(list))

	for i:=0;i<len(list);i++ {
		var user User

		// first acquire data by call
		// if failed , read in database
		result, err := ipc.Call(list[i], ipc.USERINFO_REQUEST,nil)
		if err != nil {
			user, err = user_tbl.Read(list[i])
			if err != nil {
				panic(fmt.Sprintf("cannot read user:%v in database", list[i]))
			}
		} else {
			user = result.(User)
		}

		fmt.Println(user)
		out.F_items[i].F_id = user.Id
		out.F_items[i].F_name = user.Name
		out.F_items[i].F_rank = user.Score
		out.F_items[i].F_state = int32(user.State)

		t := int32(user.ProtectTime.Unix() - time.Now().Unix())
		if t >0 {
			out.F_items[i].F_protect_time = t
		} else {
			out.F_items[i].F_protect_time = 0
		}
	}

	writer := packet.Writer()
	return pack(out, writer), nil
}
