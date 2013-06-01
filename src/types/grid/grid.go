package grid

const (
	W = byte(50)
	H = 50
)

type Grid struct {
	M []uint16
}

//------------------------------------------------ Create a new grid struct
func New() *Grid {
	g := &Grid{}
	g.M = make([]uint16, int(W)*int(H))
	return g
}

//------------------------------------------------ Set v->[X,Y]
func (g *Grid) Set(X, Y byte, v uint16) {
	if X >= 0 && X <= W {
		if Y >= 0 && Y <= H {
			g.M[Y*W+X] = v
		}
	}
}

//------------------------------------------------ Get <-[X,Y]
func (g *Grid) Get(X, Y byte) uint16 {
	if X >= 0 && X <= W {
		if Y >= 0 && Y <= H {
			return g.M[Y*W+X]
		}
	}

	return 0
}
