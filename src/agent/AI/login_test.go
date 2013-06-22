package AI

import (
	"db/data_tbl"
	"db/user_tbl"
	"testing"
	. "types"
	"types/estates"
	"types/samples"
)

func TestLoginProc(t *testing.T) {
	sess := &Session{}
	sess.MQ = make(chan IPCObject, 100)
	sess.User = &User{Id: 1}
	LoginProc(sess)

	user_tbl.Set(sess.User)
	data_tbl.Set(estates.COLLECTION, sess.Estates)
	data_tbl.Set(samples.COLLECTION, sess.LatencySamples)
}

func BenchmarkLoginProc(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		sess := &Session{}
		sess.MQ = make(chan IPCObject, 10)
		sess.User = &User{Id: int32(i)}
		LoginProc(sess)
		user_tbl.Set(sess.User)
		data_tbl.Set(estates.COLLECTION, sess.Estates)
		data_tbl.Set(samples.COLLECTION, sess.LatencySamples)
	}
}
