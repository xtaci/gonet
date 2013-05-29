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
	gaussian := Dist{}
	for i := 0; i < 1000; i++ {
		v := int16(gen.Int31n(10))
		gaussian.Add(v)
	}

	fmt.Println("N:", gaussian.n, "SIGMA:", gaussian.sigma)
	fmt.Println("Samples:", gaussian.samples)

	// testing
	fmt.Println("range [0,10]")
	for i := 0; i < 10; i++ {
		v := int16(gen.Int31n(10))
		fmt.Println(v, ":", gaussian.P(v))
	}

	fmt.Println("range [0,20]")
	for i := 0; i < 10; i++ {
		v := int16(gen.Int31n(20))
		fmt.Println(v, ":", gaussian.P(v))
	}
}
