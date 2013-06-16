package helper

import (
	"sync/atomic"
	"time"
)

var _X0 uint32 = uint32(time.Now().UnixNano())
var _a uint32 = 1664525
var _c uint32 = 1013904223

//------------------------------------------------ Linear congruential generator, very fast!!
func LCG() uint32 {
	for {
		X0 := atomic.LoadUint32(&_X0)
		Xn := _a*X0 + _c
		if atomic.CompareAndSwapUint32(&_X0, X0, Xn) {
			return Xn
		}
	}
}
