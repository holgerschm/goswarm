package goswarm

type terminator interface {
	shouldTerminate(currentBest *Candidate) bool
}
