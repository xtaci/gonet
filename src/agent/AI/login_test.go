package AI

import (
	"testing"
	. "types"
)

func TestLoginProc(t *testing.T) {
	sess := &Session{}
	sess.MQ = make(chan IPCObject, 100)
	sess.User = User{Id: 1}
	LoginProc(sess)
}
