package AI

import (
	"fmt"
	"testing"
	"time"
	. "types"
)

func _recv(CH chan IPCObject) {
	for {
		v, ok := <-CH

		if !ok {
			break
		}
		fmt.Println(v)
	}
}

func TestLoginProc(t *testing.T) {
	sess := &Session{}
	sess.MQ = make(chan IPCObject, 100)
	go _recv(sess.MQ)
	sess.User = User{Id: 1}
	LoginProc(sess)
	time.Sleep(time.Second)
	close(sess.MQ)
}
