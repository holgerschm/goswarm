package goswarm

type multiplexer struct {
	outputs []chan<- *candidate
}

func (m multiplexer) send(best *candidate) {
	for i := 0; i < len(m.outputs); i++ {
		m.outputs[i] <- best
	}
}
