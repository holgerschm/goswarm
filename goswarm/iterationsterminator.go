package goswarm

type iterationsTerminator struct {
	minIterations int64
}

func (i iterationsTerminator) shouldTerminate(currentBest *candidate) bool {
	return currentBest.iteration >= i.minIterations
}

var _ terminator = &iterationsTerminator{}
