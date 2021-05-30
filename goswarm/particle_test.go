package goswarm

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatItStartsAtARandomPosition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objective := NewFakeObjective()
	candidateInput := make(chan *candidate)
	candidateOutput := make(chan *candidate)
	resultOutput := make(chan *candidate)
	rng := NewMockrandom(ctrl)
	rng.EXPECT().next(0.0, 1.0).Return(1.1)
	rng.EXPECT().next(-1.0, 3.0).Return(2.2)
	rng.EXPECT().next(-2.0, 5.0).Return(3.3)
	rng.EXPECT().next(-3.0, 7.0).Return(4.4)

	sut := particle{objective, candidateInput, candidateOutput, resultOutput, rng}
	sut.start()

	parameters := <-objective.Parameters
	objective.Result <- 0

	assert.Equal(t, 1.1, parameters[0])
	assert.Equal(t, 2.2, parameters[1])
	assert.Equal(t, 3.3, parameters[2])
	assert.Equal(t, 4.4, parameters[3])
}

type fakeObjective struct {
	Parameters chan []float64
	Result     chan float64
}

func NewFakeObjective() *fakeObjective {
	f := fakeObjective{make(chan []float64), make(chan float64)}
	return &f
}

func (f *fakeObjective) Dimensions() int {
	return 4
}

func (f *fakeObjective) GetLowerBound(dimension int) float64 {
	return -float64(dimension)
}

func (f *fakeObjective) GetUpperBound(dimension int) float64 {
	return float64(1 + 2*dimension)
}

func (f *fakeObjective) Evaluate(parameter []float64) float64 {
	f.Parameters <- parameter
	return <-f.Result
}

var _ Objective = &fakeObjective{}
