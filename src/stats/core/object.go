package core

import (
	"encoding/json"
)

const (
	UNKNOWN = iota
	// 基本
	TYPE_LOGIN  // 登陆
	TYPE_LOGOUT // 登出
	TYPE_CHAT   // 发送一次聊天
	TYPE_EMAIL  // 发送一次邮件

	TYPE_PVP // 发生了一次PVP
	TYPE_PVE // 发生了一次PVE

	// 军队
	TYPE_TRAIN_SOLIDER // 训练完成一个小兵, Value为(兵种:数量）
	TYPE_TRAIN_HERO    // 训练完成一个英雄, Value为(兵种:数量）
	TYPE_HERO_FIGHT    //英雄出战一次

	// 生产
	TYPE_PRODUCT_NATURAL // 自然产出
	TYPE_PRODUCE_PVE     // 一次PVE 资源产出
	TYPE_PRODUCE_PVP     // 一次PVP资源产出

	// 消耗
	TYPE_CONSUME_UPGRADE // 一次升级消耗
	TYPE_CONSUME_TRAIN   // 一次训练消耗
	TYPE_CONSUME_EQUIP   //  一次强化装备消耗
	TYPE_CONSUME_SEARCH  //  一次搜寻对手消耗
	TYPE_LOST_RESOURCE   // 一次资源损失

	// 宝石
	TYPE_CONSUME_GEM // 一次宝石消耗
)

type StatsObject struct {
	UserId    int32
	Type      int32
	Property  map[string]string
	Timestamp int64
	// TODO: add fields
}

type Summary struct {
	UserId int32
	// TODO: summary fields
}

func (sum *Summary) Marshal() []byte {
	json_val, _ := json.Marshal(sum)
	return json_val
}
