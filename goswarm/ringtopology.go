package goswarm

var _ topology = &ringTopology{}

type ringTopology struct {
	particles int
}

func (f *ringTopology) particleCount() int {
	return f.particles
}

func (f *ringTopology) getOutputs(input int) []int {
	res := make([]int, 2)
	res[0] = (input + f.particles - 1) % f.particles
	res[1] = (input + 1) % f.particles
	return res
}

func newRingTopology(particleCount int) topology {
	return &ringTopology{particleCount}
}
