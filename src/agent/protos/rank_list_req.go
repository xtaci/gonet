package protos

import (
	"time"
)

import (
	. "types"
	"misc/packet"
	"agent/ipc"
	"hub/names"
	"hub/ranklist"
	"fmt"
	"errors"
	"db/user_tbl"
)

func _rank_list_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	list := ranklist.GetRankList(1, ranklist.Count())
	out := rank_list{}
	out.items = make([]rank_list_item, len(list))

	for i:=0;i<len(list);i++ {
		var user User

		// first in memory
		if peer := names.Query(list[i]);peer!=nil {
			user, err = _get_peer_info(peer)
		}

		// then in db
		if err != nil {
			user, err = user_tbl.Read(list[i])
		}

		if err != nil {
			panic(err)
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

	return
}

func _get_peer_info(peer chan interface{})(user User, err error){
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("chan temporarily failed")
		}
	}()

	req := &ipc.RequestType{Code:ipc.USERINFO_REQUEST}
	req.CH = make(chan interface{})
	peer <- req
	return (<-req.CH).(User), nil
}
