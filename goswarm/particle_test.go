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
	// start position
	rng.EXPECT().next(0.0, 10.0).Return(1.1)
	rng.EXPECT().next(-1.0, 12.0).Return(2.2)
	rng.EXPECT().next(-2.0, 14.0).Return(3.3)
	rng.EXPECT().next(-3.0, 16.0).Return(4.4)
	// start speed
	rng.EXPECT().next(-10.0, 10.0).Return(0.0)
	rng.EXPECT().next(-13.0, 13.0).Return(0.0)
	rng.EXPECT().next(-16.0, 16.0).Return(0.0)
	rng.EXPECT().next(-19.0, 19.0).Return(0.0)

	sut := newParticle(objective, candidateInput, candidateOutput, resultOutput, rng)
	sut.start()

	parameters := <-objective.Parameters
	objective.Result <- 7.7

	assert.Equal(t, 1.1, parameters[0])
	assert.Equal(t, 2.2, parameters[1])
	assert.Equal(t, 3.3, parameters[2])
	assert.Equal(t, 4.4, parameters[3])

	// output first value as best
	output := <-candidateOutput

	assert.Equal(t, 1.1, output.parameters[0])
	assert.Equal(t, 2.2, output.parameters[1])
	assert.Equal(t, 3.3, output.parameters[2])
	assert.Equal(t, 4.4, output.parameters[3])
	assert.Equal(t, 7.7, output.value)

	close(candidateInput)
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
	return float64(10 + 2*dimension)
}

func (f *fakeObjective) Evaluate(parameter []float64) float64 {
	f.Parameters <- parameter
	return <-f.Result
}

var _ Objective = &fakeObjective{}
