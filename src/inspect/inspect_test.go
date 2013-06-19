package inspect

import (
	"agent/ipc"
	"os"
	"testing"
	. "types"
)

func TestInspect(t *testing.T) {
	var sess Session
	sess.User = &User{Id: 1}
	ipc.RegisterOnline(&sess, sess.User.Id)
	Inspect(1, os.Stdout)
	InspectField(1, ".User", os.Stdout)
}
