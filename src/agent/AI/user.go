package AI

import (
	"agent/gsdb"
	"misc/geoip"
	. "types"
)

//------------------------------------------------ 登陆后的数据加载
func LoginProc(sess *Session) bool {
	if sess.User == nil {
		return false
	}

	// TODO: init data structure for session, such as builds

	// set countrycode
	if sess.User.CountryCode == "" {
		sess.User.CountryCode = geoip.Query(sess.IP)
	}

	// register as online
	gsdb.RegisterOnline(sess, sess.User.Id)

	// load messages
	go LoadIPCObjects(sess.User.Id, sess.MQ)

	return true
}
