package goswarm

var _ topology = &fullyConnectedTopology{}

type fullyConnectedTopology struct {
	particles int
}

func (f *fullyConnectedTopology) particleCount() int {
	return f.particles
}

func (f *fullyConnectedTopology) getOutputs(input int) []int {
	res := make([]int, f.particles-1)
	for i := 0; i < f.particles-1; i++ {
		if i >= input {
			res[i] = i + 1
		} else {
			res[i] = i
		}
	}
	return res
}

func newFullyConnectedTopology(particleCount int) topology {
	return &fullyConnectedTopology{particleCount}
}
