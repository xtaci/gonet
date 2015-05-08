package cfg

import "io"

type Repeater struct {
	out1, out2 io.Writer // 1-input -> 2-output
}

//----------------------------------------------- output replicator
func (r *Repeater) Write(p []byte) (int, error) {
	var n int
	var e error

	if r.out1 != nil {
		n, e = r.out1.Write(p)
	}

	if r.out2 != nil {
		n, e = r.out2.Write(p)
	}

	return n, e
}
