package core

import (
	"encoding/json"
	"sync"
)

import (
	. "types"
	"types/estates"
)

const (
	UNKNOWN = iota
	// 基本
	TYPE_LOGIN  // 登陆
	TYPE_LOGOUT // 登出
	TYPE_CHAT   // 发送一次聊天
	TYPE_EMAIL  // 发送一次邮件

	// 战斗
	TYPE_PVP        // 发生了一次PVP
	TYPE_PVE        // 发生了一次PVE
	TYPE_HERO_FIGHT //英雄出战一次

	// 军队
	TYPE_TRAIN // 训练完成一个军队单位, Property 为(兵种:数量）

	// 生产
	TYPE_PRODUCT     // 确认一次产出
	TYPE_SYS_PRODUCT // 确认一次系统产出

	// 一次资源损失
	TYPE_LOST_RESOURCE

	// GEM消耗
	TYPE_CONSUME_GEM
)

type StatsObject struct {
	Type      int32
	Property  map[string]float32
	Timestamp int64
}

//------------------------------------------------ 一个玩家对应一个
type Collector struct {
	_stats []*StatsObject
	_lock  sync.Mutex
}

func (r *Collector) Lock() {
	r._lock.Lock()
}

func (r *Collector) Unlock() {
	r._lock.Unlock()
}

type Archive struct {
	UserId    int32
	Timestamp int64
	Fields    map[string]float32
	User      User
	Estates   estates.Manager
}

func (archive *Archive) Marshal() []byte {
	json_val, _ := json.Marshal(archive)
	return json_val
}
