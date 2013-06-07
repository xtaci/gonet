package AI

import (
	"time"
)

//------------------------------------------------ 持续生产的当前量
// Sustained Production
// f(T) = Max(base(T0) + rate*(T-T0), MaxValue)
func SP(base, max int, t0 int64, rate float32) int {
	t := time.Now().Unix()
	increment := int(float32(t-t0) * rate)

	quantity := increment + base
	if quantity > max {
		quantity = max
	}

	return quantity
}
