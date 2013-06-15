package AI

import (
	"log"
	"time"
)

import (
	"agent/ipc"
	"db/data_tbl"
	. "types"
	"types/estates"
)

//------------------------------------------------ 登陆后的数据加载
func LoginProc(sess *Session) bool {
	// 载入建筑表
	if !data_tbl.Get(estates.COLLECTION, sess.User.Id, &sess.Estates) {
		log.Println("Cannot Get Estates from Database.")
	} else {
		// 创建Grid
		sess.Grid = CreateGrid(&sess.Estates)
	}

	// 开始计算Flush时间
	sess.LastFlushTime = time.Now().Unix()

	// 最后, 载入离线消息，并push到MQ, 这里小心MQ的buffer长度
	LoadIPCObjects(sess.User.Id, sess.MQ)

	// 注册为在线
	ipc.RegisterOnline(sess, sess.User.Id)
	return true
}
