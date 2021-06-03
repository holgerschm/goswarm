package goswarm

import (
	"math/rand"
	"time"
)

type randwrapper struct {
	rng *rand.Rand
}

func (s *randwrapper) next(lower float64, upper float64) float64 {
	return s.rng.Float64()*(upper-lower) + lower
}

func NewSystemRandom() random {
	rw := randwrapper{}
	time.Sleep(time.Millisecond)
	rw.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &rw
}

var _ random = &randwrapper{}
