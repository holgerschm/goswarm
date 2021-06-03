package goswarm

type multiplexer interface {
	send(best *Candidate)
}

type nonBlockingMultiplexer struct {
	outputs []chan<- *Candidate
}

func (m *nonBlockingMultiplexer) send(best *Candidate) {
	for i := 0; i < len(m.outputs); i++ {
		select {
		case m.outputs[i] <- best:
		default:
		}
	}
}
