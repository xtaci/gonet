package AI

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	//"fmt"
)

import (
	"db/data_tbl"
	"db/forward_tbl"
	"gamedata"
	. "types"
	"types/estates"
	"types/grid"
)

//------------------------------------------------ 登陆后的数据加载
func LoginWork(sess *Session) bool {
	// 载入建筑表
	data_tbl.Get(estates.COLLECTION, sess.User.Id, &sess.Estates)
	// 建立位图的格子信息
	sess.Grid = grid.New()
	for k, v := range sess.Estates.Estates {
		// TODO :  读gamedata,建立grid信息
		name := gamedata.Query(v.TYPE)
		cell := gamedata.GetString("建筑规格", name, "占用格子数")
		wh := strings.Split(cell, "X")
		w, _ := strconv.Atoi(wh[0])
		h, _ := strconv.Atoi(wh[1])

		oid, _ := strconv.Atoi(k)
		//	fmt.Println(w,h, wh, cell, v)
		for x := v.X; x < v.X+byte(w); x++ {
			for y := v.Y; y < v.Y+byte(h); y++ {
				sess.Grid.Set(x, y, uint16(oid))
			}
		}
	}

	// 最后, 载入离线消息，并push到MQ, 这里小心MQ的buffer长度
	objs := forward_tbl.PopAll(sess.User.Id)
	for k := range objs {
		obj := &IPCObject{}
		err := json.Unmarshal(objs[k], obj)
		if err != nil {
			log.Println("Illegal IPCObject", objs[k])
		} else {
			sess.MQ <- *obj
		}
	}

	sess.LastFlushTime = time.Now().Unix()

	return true
}
