package AI

import (
	"helper"
	"math"
	"time"
)

//------------------------------------------------ 持续生产的当前量
// Sustained Production
// f(T) = Min(base(T0) + rate*(T-T0), MaxValue)
func SP(base, max int, t0 int64, rate float32) int {
	t := time.Now().Unix()
	increment := int(float32(t-t0) * rate)

	quantity := increment + base
	if quantity > max {
		quantity = max
	}

	return quantity
}

//------------------------------------------------ 随机概率丢骰子
func Dice(probability float32) bool {
	return helper.LCG() < uint32(probability*math.MaxUint32)
}
