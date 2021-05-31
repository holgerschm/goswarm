package goswarm

var _ topology = &fullyConnectedTopology{}

type fullyConnectedTopology struct {
	particles int
}

func (f *fullyConnectedTopology) particleCount() int {
	return f.particles
}

func (f *fullyConnectedTopology) getOutputs(input int) []int {
	res := make([]int, f.particles)
	for i := 0; i < f.particles; i++ {
		res[i] = i
	}
	return res
}

func newFullyConnectedTopology(particleCount int) topology {
	return &fullyConnectedTopology{particleCount}
}
