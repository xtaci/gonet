package grid

const (
	W = 50
	H = 50
)

type Grid struct {
	_m []uint16
}

//------------------------------------------------ Create a new grid struct
func New() *Grid {
	g := &Grid{}
	g._m = make([]uint16, W*H)
	return g
}

//------------------------------------------------ Set v->[X,Y]
func (g *Grid) Set(X, Y int, v uint16) {
	if X >= 0 && X < W {
		if Y >= 0 && Y < H {
			g._m[Y*W+X] = v
		}
	}
}

//------------------------------------------------ Get <-[X,Y]
func (g *Grid) Get(X, Y int) uint16 {
	if X >= 0 && X < W {
		if Y >= 0 && Y < H {
			return g._m[Y*W+X]
		}
	}

	return 0
}
