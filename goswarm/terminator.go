package goswarm

type terminator interface {
	shouldTerminate(currentBest *candidate) bool
}
