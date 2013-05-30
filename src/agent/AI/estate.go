package AI

import (
	"time"
)

import (
	event "agent/event_client"
	"db/estate_tbl"
	"gamedata"
	. "types"
	"types/estate"
)

//------------------------------------------------ 创建新的建筑
func EstateNew(sess *Session, name string, X, Y uint16) bool {
	// 获取资源消耗，检查当前资源是否满足一级建筑的建造条件
	fields := gamedata.FieldNames(ESTATE_TBL)
	for _, v := range fields {
		cost := gamedata.GetInt(ESTATE_TBL, 1, v)
		if c := sess.Res.Get(v); c < cost {
			return false
		} else { // 扣除资源
			sess.Res.Set(v, c-cost)
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
	E.Timeout = time.Now().Unix() + cd_time
	event_id := event.Add(N.OID, sess.Basic.Id, E.Timeout)

	// 变更当前session内容
	sess.EstateManager.AppendEstate(N)
	sess.EstateManager.AppendCD(event_id, E)

	// 持久化
	return estate_tbl.Set(&sess.EstateManager)
}
