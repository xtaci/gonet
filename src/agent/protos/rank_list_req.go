package protos

import (
	. "types"
	"misc/packet"
	"agent/ipc"
	"hub/names"
	"hub/ranklist"
	"fmt"
	"errors"
)

func _rank_list_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	list := ranklist.GetRankList(1, ranklist.Count())
	tbl := rank_list{}
	tbl.items = make([]rank_list_item, len(list))


	for i:=0;i<len(list);i++ {
		var user User
		var err error

		// first in memory
		if peer := names.Query(list[i]);peer!=nil {
			user, err = _get_peer_info(peer)
		}

		// then in db
		if err != nil {
			
		}

		fmt.Println(user)
		/*
		tbl.items[i].id = list[i]
		tbl.items[i].name =
		tbl.items[i].rank =
		tbl.items[i].state =
		tbl.items[i].protect_time =
		*/
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
