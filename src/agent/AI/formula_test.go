package AI

import (
	"fmt"
	"testing"
	"time"
)

func TestSP(t *testing.T) {
	t0 := time.Now().Unix()

	for i := int64(0); i < 150; i++ {
		fmt.Println(SP(0, 10000, t0-i, 100))
	}
}
