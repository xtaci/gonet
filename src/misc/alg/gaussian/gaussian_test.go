package gaussian

import (
	"testing"
	"fmt"
)

func TestGauss(t *testing.T) {
	for i:=-5.0;i<=5.0;i+=0.1 {
		fmt.Println(StdDeviation(float64(i)))
	}
}
