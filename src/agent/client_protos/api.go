package protos

import "misc/packet"
import . "types"

var Code map[string]int16 = map[string]int16 {
	"heart_beat_req":0,	// 心跳包..
	"user_login_req":1,	// 客户端发送用户登陆请求包
	"user_login_succeed_ack":2,	// 登陆成功
	"user_login_faild_ack":3,	// 登陆失败
	"save_user_info_req":4,	// 存档
	"rank_list_req":5,	// 客户端向服务器请求排行榜快照.
	"rank_list_ack":6,	// 排行榜信息.
	"pve_list_req":7,	// 客户端向服务器请求PVE快照.
	"pve_list_ack":8,	// pve信息
	"atk_player_req":9,	// 攻击另一玩家
	"atk_player_succeed_ack":10,	// 攻击玩家成功
	"atk_player_faild_ack":11,	// 攻击玩家失败
	"atk_player_rst_req":12,	// 攻击玩家结果存档
	"atk_monster_req":13,	// 攻击怪物
	"atk_monster_ack":14,	// 攻击怪物结果
	"atk_monster_rst_req":15,	// 攻击怪物结果存档
}

var RCode map[int16]string = map[int16]string {
	0:"heart_beat_req",	// 心跳包..
	1:"user_login_req",	// 客户端发送用户登陆请求包
	2:"user_login_succeed_ack",	// 登陆成功
	3:"user_login_faild_ack",	// 登陆失败
	4:"save_user_info_req",	// 存档
	5:"rank_list_req",	// 客户端向服务器请求排行榜快照.
	6:"rank_list_ack",	// 排行榜信息.
	7:"pve_list_req",	// 客户端向服务器请求PVE快照.
	8:"pve_list_ack",	// pve信息
	9:"atk_player_req",	// 攻击另一玩家
	10:"atk_player_succeed_ack",	// 攻击玩家成功
	11:"atk_player_faild_ack",	// 攻击玩家失败
	12:"atk_player_rst_req",	// 攻击玩家结果存档
	13:"atk_monster_req",	// 攻击怪物
	14:"atk_monster_ack",	// 攻击怪物结果
	15:"atk_monster_rst_req",	// 攻击怪物结果存档
}

var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet)([]byte, error){
	0:P_heart_beat_req,
	1:P_user_login_req,
	4:P_save_user_info_req,
	5:P_rank_list_req,
	7:P_pve_list_req,
	9:P_atk_player_req,
	12:P_atk_player_rst_req,
	13:P_atk_monster_req,
	15:P_atk_monster_rst_req,
}