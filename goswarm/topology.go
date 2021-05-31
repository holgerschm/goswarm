package goswarm

type topology interface {
	particleCount() int
	getOutputs(input int) []int
}
