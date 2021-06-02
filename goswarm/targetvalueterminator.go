package goswarm

type targetValueTerminator struct {
	limit float64
}

func (t targetValueTerminator) shouldTerminate(currentBest *candidate) bool {
	return currentBest.value < t.limit
}

var _ terminator = &targetValueTerminator{}
