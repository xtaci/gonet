package gaussian

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestGauss(t *testing.T) {
	src := rand.NewSource(time.Now().Unix())
	gen := rand.New(src)

	// gaussian
	gaussian := NewDist(128)
	for i := 0; i < 1000; i++ {
		v := gen.Intn(200)
		gaussian.Add(v)
	}

	fmt.Println("N-samples:", gaussian.N, ", σ:", gaussian.Sigma)

	// testing
	fmt.Println("range [0,200]")
	sigma := gaussian.Sigma
	mean := gaussian.Mean
	for i := 0; i < 10; i++ {
		v := gen.Intn(200)
		fmt.Printf("X:%4d: P(v)=%0.4f, deriv:%.2fσ\n", v, gaussian.P(v), math.Abs(float64(v)-mean)/sigma)
	}

	fmt.Println("range [0,400]")
	for i := 0; i < 10; i++ {
		v := gen.Intn(400)
		fmt.Printf("X:%4d: P(v)=%0.4f, deriv:%.2fσ\n", v, gaussian.P(v), math.Abs(float64(v)-mean)/sigma)
	}
}
