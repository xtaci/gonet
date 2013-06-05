package protos

const (
	TYPE_PVP                     = iota // 发生了一次PVP
	TYPE_PVE                            // 发生了一次PVE
	TYPE_TRAIN                          // 训练了一个小兵
	TYPE_HERO_ATTACk                    // 英雄出战一次
	TYPE_CONSUME_FOOD                   // 消耗一次食物
	TYPE_CONSUME_GOLD                   // 消耗一次金币
	TYPE_CONSUME_GOLD_AS_ENHANCE        // 消耗一次强化
	TYPE_ROB_FOOD                       // 掠夺一次的食物量
	TYPE_ROB_GOLD                       // 掠夺一次的金币量
)

type StatsObject struct {
	UserId int32
	Type   int32
	// TODO: add fields
}
