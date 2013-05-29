package gaussian

import (
	"fmt"
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
		v := int16(gen.Int31n(200))
		fmt.Println(v)
		gaussian.Add(v)
	}

	fmt.Println("N:", gaussian.n, "SIGMA:", gaussian.sigma)
	fmt.Println("Samples:", gaussian.samples)

	// testing
	fmt.Println("range [0,200]")
	for i := 0; i < 10; i++ {
		v := int16(gen.Int31n(200))
		fmt.Println(v, ":", gaussian.P(v))
	}

	fmt.Println("range [0,1000]")
	for i := 0; i < 10; i++ {
		v := int16(gen.Int31n(1000))
		fmt.Println(v, ":", gaussian.P(v))
	}
}
