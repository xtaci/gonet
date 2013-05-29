package gaussian

import "math"

const (
	SQRT2PI = float64(2.506628274631001)
	SIGMA = float64(1.0)
	MU = float64(0.0)
	MAX_SAMPLES = 128
)

type Dist struct {
	samples [MAX_SAMPLES]int64
	ptr int
	n	int
	sigma float64
}

func (dist *Dist) Add(x int64) {
	dist.samples[dist.ptr] = x
	if dist.ptr++;dist.ptr >= MAX_SAMPLES {
		dist.ptr = 0
	}

	if dist.n < MAX_SAMPLES {
		dist.n++
	}

	// caculate mean
	sum := int64(0)
	for i:=0;i<dist.n;i++ {
		sum += dist.samples[i]
	}

	mean := float64(sum)/float64(dist.n)

	// caculate standard deviation 
	sum2 := float64(0.0)
	for i:=0;i<dist.n;i++ {
		v := float64(dist.samples[i]) - mean
		v = v*v
		sum2 += v
	}

	dist.sigma = math.Sqrt(sum2/float64(dist.n))
}

func (dist *Dist) P(x int64) float64 {
	X := float64(x)
	A := 1.0/(dist.sigma * SQRT2PI)
	B := math.Exp(-(X*X)/(2*dist.sigma*dist.sigma))
	return A*B
}
