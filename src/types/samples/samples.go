package samples

import (
	"strconv"
)

import (
	"cfg"
	"misc/alg/gaussian"
)

const COLLECTION = "SAMPLES"

type Samples struct {
	UserId  int32
	Version uint32
	G       *gaussian.Dist
}

func (s *Samples) Init() {
	config := cfg.Get()
	samples, _ := strconv.Atoi(config["samples"])

	s.G = gaussian.NewDist(samples)
}
