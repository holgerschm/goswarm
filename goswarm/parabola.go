package goswarm

// for testing purposes -> minimum at f(2) = -1
type parabola struct {
}

func (p parabola) Dimensions() int {
	return 1
}

func (p parabola) GetLowerBound(dimension int) float64 {
	return -1000
}

func (p parabola) GetUpperBound(dimension int) float64 {
	return 5000
}

func (p parabola) Evaluate(parameter []float64) float64 {
	x := parameter[0]

	return (x - 1) * (x - 3)
}
