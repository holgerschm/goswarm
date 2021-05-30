package goswarm

type Objective interface {
	Dimensions() int
	GetLowerBound(dimension int) float64
	GetUpperBound(dimension int) float64
	Evaluate(parameter []float64) float64
}
