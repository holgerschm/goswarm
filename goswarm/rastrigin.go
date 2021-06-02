package goswarm

import "math"

// for testing purposes -> minimum at f(0,0,..,0) = 0
type rastrigin struct {
}

func (p rastrigin) Dimensions() int {
	return 13
}

func (p rastrigin) GetLowerBound(dimension int) float64 {
	return -5.12
}

func (p rastrigin) GetUpperBound(dimension int) float64 {
	return 5.12
}

func (p rastrigin) Evaluate(parameter []float64) float64 {
	n := p.Dimensions()

	result := float64(10 * n)
	for i := 0; i < n; i++ {
		x := parameter[i]
		result += x * x
		result -= 10 * math.Cos(2*math.Pi*x)
	}

	return result
}
