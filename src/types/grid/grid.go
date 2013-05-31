package grid

const (
	W = byte(50)
	H = 50
)

type Grid struct {
	M []uint32
}

//------------------------------------------------ Create a new grid struct
func New() *Grid {
	g := &Grid{}
	g.M= make([]uint32, int(W)*int(H))
	return g
}

//------------------------------------------------ Set v->[X,Y]
func (g *Grid) Set(X,Y byte, v uint32) {
	g.M[Y*W + X]= v
}

//------------------------------------------------ Get <-[X,Y]
func (g *Grid) Get(X,Y byte) uint32 {
	return g.M[Y*W + X]
}
