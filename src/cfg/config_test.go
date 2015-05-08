package cfg

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := Get()

	for k, v := range config {
		fmt.Println(k, "=", v)
	}
}
