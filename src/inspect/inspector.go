package inspect

import (
	"encoding/json"
	"fmt"
	"io"
)

import (
	"agent/ipc"
)

func Inspect(id int32, output io.Writer) {
	sess := ipc.QueryOnline(id)
	if sess != nil {
		val, _ := json.Marshal(sess.User)
		fmt.Fprintln(output, string(val))
		val, _ = json.Marshal(sess.Estates)
		fmt.Fprintln(output, string(val))
		val, _ = json.Marshal(sess.Soldiers)
		fmt.Fprintln(output, string(val))
	}
}

func ListAll(output io.Writer) {
	fmt.Fprintln(output, ipc.ListAll())
}
