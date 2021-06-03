package goswarm

type iterationsTerminator struct {
	minIterations int64
}

func (i iterationsTerminator) shouldTerminate(currentBest *Candidate) bool {
	return currentBest.Iteration >= i.minIterations
}

var _ terminator = &iterationsTerminator{}
