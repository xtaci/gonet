package protos

import "misc/packet"

import . "types"

var Code map[string]uint16 = map[string]uint16{
	"heart_beat_req":         0,  // payload:null 心跳包..
	"user_login_req":         1,  // payload:user_login_info 客户端发送用户登陆请求包
	"user_login_succeed_ack": 2,  // payload:user_snapshot 登陆成功
	"user_login_faild_ack":   3,  // payload:command_result_pack 登陆失败
	"save_user_info_req":     4,  // payload:user_archives_info 存档
	"rank_list_req":          5,  // payload:null 客户端向服务器请求排行榜快照.
	"rank_list_ack":          6,  // payload:rank_list 排行榜信息.
	"pve_list_req":           7,  // payload:null 客户端向服务器请求PVE快照.
	"pve_list_ack":           8,  // payload:pve_list pve信息
	"atk_player_req":         9,  // payload:command_id_pack 攻击另一玩家
	"atk_player_succeed_ack": 10, // payload:user_snapshot 攻击玩家成功
	"atk_player_faild_ack":   11, // payload:command_result_pack 攻击玩家失败
	"atk_player_rst_req":     12, // payload:atk_player_rst_req 攻击玩家结果存档
	"atk_monster_req":        13, // payload:command_id_pack 攻击怪物
	"atk_monster_ack":        14, // payload:command_result_pack 攻击怪物结果
	"atk_monster_rst_req":    15, // payload:atk_monster_rst_req 攻击怪物结果存档
}

var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error) = map[uint16]func(*Session, *packet.Packet) ([]byte, error){
	0:  _heart_beat_req,
	1:  _user_login_req,
	4:  _save_user_info_req,
	5:  _rank_list_req,
	7:  _pve_list_req,
	9:  _atk_player_req,
	12: _atk_player_rst_req,
	13: _atk_monster_req,
	15: _atk_monster_rst_req,
}
