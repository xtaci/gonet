package grid

import (
	"fmt"
	"testing"
)

func TestGrid(t *testing.T) {
	grid := New()
	for y := 0; y < 60; y++ {
		for x := 0; x < 60; x++ {
			grid.Set(x, y, uint16(x))
		}
	}

	for y := 0; y < 60; y++ {
		for x := 0; x < 60; x++ {
			fmt.Print(grid.Get(x, y), " ")
		}
		fmt.Println()
	}

	fmt.Println()

	fmt.Print(grid._m)
}
