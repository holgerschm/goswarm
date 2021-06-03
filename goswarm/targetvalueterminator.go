package goswarm

type targetValueTerminator struct {
	limit float64
}

func (t targetValueTerminator) shouldTerminate(currentBest *Candidate) bool {
	return currentBest.Value < t.limit
}

var _ terminator = &targetValueTerminator{}
