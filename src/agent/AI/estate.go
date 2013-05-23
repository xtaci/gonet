package AI

import (
	"time"
)

import (
	. "types"
	"types/estate"
	"gamedata"
	event "agent/event_client"
	"db/estate_tbl"
)

//------------------------------------------------ 创建新的建筑
func EstateNew(sess *Session, name string, X,Y uint16) bool {
	// 获取资源消耗，检查当前资源是否满足一级建筑的建造条件
	fields := gamedata.FieldNames(ESTATE_TBL)
	for _,v := range fields {
		cost := gamedata.GetInt(ESTATE_TBL,1,v)
		if sess.Res.Get(v) < cost {
			return false
		}
	}

	// TODO: 解锁检查
	// TODO: 冷却时间表读取
	// TODO: 检查是否放得下 
	cd_time := int64(10)
	// 创建新的建筑
	N := &estate.Estate{}
	N.OID = sess.EstateManager.GENID()
	N.Status = estate.STATUS_CD
	N.X = X
	N.Y = Y

	// 新的冷却事件
	E := &estate.CD{}
	E.OID = N.OID
	E.EventId = event.Add(N.OID, sess.Basic.Id, time.Now().Unix()+cd_time)

	// 变更当前session内容
	sess.EstateManager.AppendEstate(N)
	sess.EstateManager.AppendCD(E)

	// 持久化
	return estate_tbl.Set(sess.Basic.Id, sess.EstateManager)
}
