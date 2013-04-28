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
	out.items = make([]rank_list_item, len(list))

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
		out.items[i].id = user.Id
		out.items[i].name = user.Name
		out.items[i].rank = user.Score
		out.items[i].state = int32(user.State)

		t := int32(user.ProtectTime.Unix() - time.Now().Unix())
		if t >0 {
			out.items[i].protect_time = t
		} else {
			out.items[i].protect_time = 0
		}
	}

	writer := packet.Writer()
	return pack(out, writer), nil
}
