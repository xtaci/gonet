package main

import "io"

type Repeater struct {
	Out1, Out2 io.Writer
}

func (r *Repeater) Write(p []byte) (n int, err error) {
	r.Out1.Write(p)
	n, err = r.Out2.Write(p)

	return
}
