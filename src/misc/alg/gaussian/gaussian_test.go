package gaussian

import (
	"testing"
	"fmt"
	"math/rand"
	"time"
)

func TestGauss(t *testing.T) {
	src := rand.NewSource(time.Now().Unix())
	gen := rand.New(src)

	// gaussian
	gaussian := Dist{}
	for i:=0;i<1000;i++ {
		v := gen.Int63n(10)
		gaussian.Add(v)
	}

	fmt.Println("N:", gaussian.n, "SIGMA:", gaussian.sigma)
	fmt.Println("Samples:", gaussian.samples)
	sum := int64(0)
	for _,v := range gaussian.samples {
		sum+=v
	}

	// testing
	fmt.Println("range [0,10]")
	for i:=0;i<10;i++ {
		v := gen.Int63n(10)
		fmt.Println(v, ":", gaussian.P(v))
	}

	fmt.Println("range [0,20]")
	for i:=0;i<10;i++ {
		v := gen.Int63n(20)
		fmt.Println(v, ":", gaussian.P(v))
	}
}
