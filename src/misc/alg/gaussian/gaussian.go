package gaussian

import "math"

const (
	SQRT2PI = float64(2.506628274631001)
	SIGMA = float64(1.0)
	MU = float64(0.0)
)

//------------------------------------------------ Standard Deviation
func StdDeviation(x float64) float64 {
	return P(x, MU, SIGMA)
}

func P(x float64, mu float64, sigma float64) float64 {
	A := 1.0/(sigma * SQRT2PI)
	B := math.Exp(-((x-mu)*(x-mu))/(2*sigma*sigma))
	return A*B
}
