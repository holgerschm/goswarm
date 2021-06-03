package goswarm

import "fmt"

type Logger interface {
	Log(currentBest Candidate)
}

type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(currentBest Candidate) {
	fmt.Println("iteration:", currentBest.Iteration)
	fmt.Println("parameters:", currentBest.Parameters)
	fmt.Println("value:", currentBest.Value)
}

var _ Logger = &ConsoleLogger{}

type nilLogger struct{}

func (n *nilLogger) Log(currentBest Candidate) {
}

var _ Logger = &nilLogger{}
