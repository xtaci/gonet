package inspect

import (
	"agent/ipc"
	"io"
	"fmt"
)

func Inspect(id int32, output io.Writer) {
	if sess:=ipc.QueryOnline(id); sess!=nil {
		fmt.Fprintln(output, sess)
	} else {
		fmt.Fprintln(output, "user is offline")
	}
}
