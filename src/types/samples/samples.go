package samples

import (
	"log"
	"strconv"
)

import (
	"cfg"
	"misc/alg/gaussian"
)

const (
	COLLECTION      = "SAMPLES"
	DEFAULT_SAMPLES = 128
)

type Samples struct {
	UserId  int32
	Version uint32
	G       *gaussian.Dist
}

func (s *Samples) Init() {
	config := cfg.Get()
	samples, err := strconv.Atoi(config["samples"])
	if err != nil {
		log.Println("cannot parse samples from config", err)
		samples = DEFAULT_SAMPLES
	}

	s.G = gaussian.NewDist(samples)
}
