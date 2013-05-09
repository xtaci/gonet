package main

import "io"

type Repeater struct {
	Out1, Out2 io.Writer // 1-input -> 2-output
}

//----------------------------------------------- output replicator
func (r *Repeater) Write(p []byte) (n int, err error) {
	r.Out1.Write(p)
	n, err = r.Out2.Write(p)

	return
}
