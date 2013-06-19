package AI

import (
	"time"
)

import (
	"agent/ipc"
	"db/data_tbl"
	"misc/geoip"
	. "types"
	"types/estates"
	"types/samples"
)

//------------------------------------------------ 登陆后的数据加载
func LoginProc(sess *Session) bool {
	if sess.User == nil {
		return false
	}

	// 载入建筑表
	if !data_tbl.Get(estates.COLLECTION, sess.User.Id, &sess.Estates) {
		// 创建默认的建筑表
		e := &estates.Manager{}
		e.UserId = sess.User.Id
		sess.Estates = e
	} else {
		// 创建Grid
		sess.Grid = CreateGrid(sess.Estates)
	}

	if !data_tbl.Get(samples.COLLECTION, sess.User.Id, &sess.LatencySamples) {
		s := &samples.Samples{}
		s.UserId = sess.User.Id
		s.Init()
		sess.LatencySamples = s
	}

	//
	if sess.User.CountryCode == "" {
		sess.User.CountryCode = geoip.Query(sess.IP)
	}

	if sess.User.CreatedAt == 0 {
		sess.User.CreatedAt = time.Now().Unix()
	}

	// 开始计算Flush时间
	sess.LastFlushTime = time.Now().Unix()

	// 最后, 载入离线消息，并push到MQ, 这里小心MQ的buffer长度
	LoadIPCObjects(sess.User.Id, sess.MQ)

	// 注册为在线
	ipc.RegisterOnline(sess, sess.User.Id)

	// TODO: 标记在线 
	//sess.LoggedIn = true
	return true
}
