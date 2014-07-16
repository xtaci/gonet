package gaussian

import "math"

const (
	SQRT2PI = float64(2.50662827463100050241576528481)
)

type Dist struct {
	Samples []int
	Ptr     int
	N       int
	Sigma   float64
	Mean    float64
}

func NewDist(num_samples int) *Dist {
	dist := &Dist{}
	dist.Samples = make([]int, num_samples)
	return dist
}

func (dist *Dist) IsSampleOk() bool {
	if dist.N >= len(dist.Samples) {
		return true
	} else {
		return false
	}
}

func (dist *Dist) Add(x int) {
	dist.Samples[dist.Ptr] = x
	if dist.Ptr++; dist.Ptr >= len(dist.Samples) {
		dist.Ptr = 0
	}

	if dist.N < len(dist.Samples) {
		dist.N++
	}

	if dist.N == len(dist.Samples) {
		// caculate mean
		sum := int64(0)
		for i := 0; i < dist.N; i++ {
			sum += int64(dist.Samples[i])
		}

		dist.Mean = float64(sum) / float64(dist.N)

		// caculate standard deviation
		sum2 := float64(0.0)
		for i := 0; i < dist.N; i++ {
			v := float64(dist.Samples[i]) - dist.Mean
			v = v * v
			sum2 += v
		}

		dist.Sigma = math.Sqrt(sum2 / float64(dist.N))
	}
}

func (dist *Dist) P(x int) float64 {
	X := float64(x)
	A := 1.0 / (dist.Sigma * SQRT2PI)
	B := math.Exp(-((X - dist.Mean) * (X - dist.Mean)) / (2 * dist.Sigma * dist.Sigma))
	return A * B
}
