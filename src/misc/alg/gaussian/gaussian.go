package gaussian

import "math"

const (
	SQRT2PI = float64(2.506628274631001)
)

type Dist struct {
	samples []int16
	ptr     int
	n       int
	sigma   float64
	mean    float64
}

func NewDist(num_samples int) *Dist {
	dist := &Dist{}
	dist.samples = make([]int16, num_samples)
	return dist
}

func (dist *Dist) Sigma() float64 {
	return dist.sigma
}

func (dist *Dist) Mean() float64 {
	return dist.mean
}

func (dist *Dist) IsSampleOk() bool {
	if dist.n >= len(dist.samples) {
		return true
	} else {
		return false
	}
}

func (dist *Dist) Add(x int16) {
	dist.samples[dist.ptr] = x
	if dist.ptr++; dist.ptr >= len(dist.samples) {
		dist.ptr = 0
	}

	if dist.n < len(dist.samples) {
		dist.n++
	}

	if dist.n == len(dist.samples) {
		// caculate mean
		sum := int64(0)
		for i := 0; i < dist.n; i++ {
			sum += int64(dist.samples[i])
		}

		dist.mean = float64(sum) / float64(dist.n)

		// caculate standard deviation
		sum2 := float64(0.0)
		for i := 0; i < dist.n; i++ {
			v := float64(dist.samples[i]) - dist.mean
			v = v * v
			sum2 += v
		}

		dist.sigma = math.Sqrt(sum2 / float64(dist.n))
	}
}

func (dist *Dist) P(x int16) float64 {
	X := float64(x)
	A := 1.0 / (dist.sigma * SQRT2PI)
	B := math.Exp(-((X - dist.mean) * (X - dist.mean)) / (2 * dist.sigma * dist.sigma))
	return A * B
}
