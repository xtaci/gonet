package cfg

import "io"

type Repeater struct {
	out1, out2 io.Writer // 1-input -> 2-output
}

//----------------------------------------------- output replicator
func (r *Repeater) Write(p []byte) (int, error) {
	r.out1.Write(p)
	return r.out2.Write(p)
}
