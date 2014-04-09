package ipc

import (
	"encoding/json"
	"time"
)

import (
	"agent/gsdb"
	"agent/hub_client"
	"db/forward_tbl"
	. "helper"
	. "types"
)

const (
	SVC_PING                        = int16(iota) // 协议ping, 测试用
	SVC_CHAT                                      // 聊天消息
	SVC_NEWMAIL                                   // 邮件
	SVC_NOTIFY12PM                                // 通知今日PM:12到达
	SVC_NOTIFYALLIANCE                            // 通知联盟信息变更
	SVC_NOTIFYALLIANCEREQUEST                     // 通知联盟信息请求
	SVC_NOTIFYLEAVELALLIANCEREQUEST               //通知退出联盟
	SVC_KICK                                      // 踢人
	SVC_MUTE                                      // 禁言
	SVC_UNMUTE                                    // 解除禁言
	SVC_BAN                                       // 禁止登录
	SVC_NEWITEM                                   // 新道具产生

	// 系统服务
	SYS_BROADCAST // 系统进程专有服务, 转发实时广播
	SYS_MULTICAST // 系统进程专有服务, 多播
	SYS_RELOADGD  // 执行GAMEDATA的重新载入
	// SYS_RELOADBANLIST  // 执行BANLIST的重新载入
	SYS_RELOADACTIVITY // 执行ACTIVITY的重新载入
	SYS_CUPCHANGE      // 同步排名变更
	SYS_CHATSYNC       // 同步频道聊天

	SVC_NEWGIFT          = 100 // 发放礼包
	SVC_NOTIFYTASKUPDATE = 101 // 通知任务数据更新
)

//---------------------------------------------------------- 异步消息发送
func Send(src_id, dest_id int32, service int16, object interface{}) (ret bool) {
	// 况序列化被传输对象为json
	val, err := json.Marshal(object)
	if err != nil {
		ERR("cannot marshal object to json", err)
		return false
	}
	// Send函数不能投递到SYS_USR
	if dest_id == SYS_USR {
		ERR("cannot Send to SYS_USR")
		return false
	}

	// 打包为IPCObject
	req := &IPCObject{SrcID: src_id,
		DestID:  dest_id,
		Service: service,
		Object:  val,
		Time:    time.Now().Unix()}

	peer := gsdb.QueryOnline(dest_id)
	if peer != nil { // 如果玩家在本服务器
		// 对方的channel 可能会close, 需要处理panic的情况
		defer func() {
			if x := recover(); x != nil {
				ret = false
				forward_tbl.Push(req)
			}
		}()
		select {
		case peer.MQ <- *req:
		case <-time.After(time.Second):
			panic("deadlock") // rare case, when both chans are full.
		}
		return true
	} else { // 通过HUB转发IPCObject
		return hub_client.Forward(req)
	}

	return false
}

//---------------------------------------------------------- 全服广播一条动态消息
func Broadcast(service int16, object interface{}) (ret bool) {
	obj, ok := _create_broadcast_ipcobject(service, object)
	if !ok {
		return false
	}

	// 投递到HUB
	hub_client.Forward(obj)

	// 投递到本地SYS_ROUTINE
	peer := gsdb.QueryOnline(SYS_USR)
	peer.MQ <- *obj
	return true
}

//---------------------------------------------------------- 本地广播一条动态消息
func Localcast(service int16, object interface{}) (ret bool) {
	obj, ok := _create_broadcast_ipcobject(service, object)
	if !ok {
		return false
	}

	// 只投递到本地SYS_ROUTINE
	peer := gsdb.QueryOnline(SYS_USR)
	peer.MQ <- *obj
	return true
}

func _create_broadcast_ipcobject(service int16, object interface{}) (*IPCObject, bool) {
	// object序列化
	val, err := json.Marshal(object)
	if err != nil {
		ERR("cannot marshal object to json", err)
		return nil, false
	}

	// 内容包
	content := &IPCObject{
		Service: service,
		Object:  val,
		Time:    time.Now().Unix(),
	}

	content_json, _ := json.Marshal(content)

	// 封装为Broadcast包
	bc := &IPCObject{
		Service: SYS_BROADCAST,
		Object:  content_json,
		Time:    time.Now().Unix(),
	}

	return bc, true
}

//---------------------------------------------------------- 组播
// 发消息到一组**给定的**目标ID
func Multicast(ids []int32, service int16, object interface{}) (ret bool) {
	// object序列化
	val, err := json.Marshal(object)
	if err != nil {
		ERR("cannot marshal object to json", err)
		return false
	}

	// 内容包
	content := &IPCObject{
		Service: service,
		Object:  val,
		Time:    time.Now().Unix(),
	}

	content_json, _ := json.Marshal(content)

	// 封装为Multicast包
	mc := &IPCObject{
		AuxIDs:  ids,
		Service: SYS_MULTICAST,
		Object:  content_json,
		Time:    time.Now().Unix(),
	}

	// 投递到HUB
	hub_client.Forward(mc)

	// 投递到本地SYS_ROUTINE
	peer := gsdb.QueryOnline(SYS_USR)
	peer.MQ <- *mc
	return true
}
