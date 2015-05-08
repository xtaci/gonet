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
	SVC_PING = int16(iota) // 协议ping, 测试用
	SVC_CHAT               // 聊天消息
	SVC_KICK               // 踢人

	// 系统服务
	SYS_BROADCAST // 系统进程专有服务, 转发实时广播
	SYS_MULTICAST // 系统进程专有服务, 多播
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
