package goswarm

type multiplexer interface {
	send(best *candidate)
}

type nonBlockingMultiplexer struct {
	outputs []chan<- *candidate
}

func (m *nonBlockingMultiplexer) send(best *candidate) {
	for i := 0; i < len(m.outputs); i++ {
		select {
		case m.outputs[i] <- best:
		default:
		}
	}
}
