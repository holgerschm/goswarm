package goswarm

import "math"

// levi no 13 for testing purposes -> minimum at f(1,1) = 0
type levi struct {
}

func (p levi) Dimensions() int {
	return 2
}

func (p levi) GetLowerBound(dimension int) float64 {
	return -10
}

func (p levi) GetUpperBound(dimension int) float64 {
	return 10
}

func (p levi) Evaluate(parameter []float64) float64 {
	x := parameter[0]
	y := parameter[1]

	s1 := math.Sin(3 * math.Pi * x)
	x1 := x - 1
	s2 := math.Sin(3 * math.Pi * y)
	y1 := y - 1
	s3 := math.Sin(2 * math.Pi * y)

	return s1*s1 + x1*x1*(1+s2*s2) + y1*y1*(1+s3*s3)
}
