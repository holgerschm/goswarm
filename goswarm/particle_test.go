package goswarm

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
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
	rng.EXPECT().next(-2.0, 138.0).Return(3.3)
	rng.EXPECT().next(-3.0, 1468.0).Return(4.4)
	// start speed
	rng.EXPECT().next(-10.0, 10.0).Return(0.0)
	rng.EXPECT().next(-13.0, 13.0).Return(0.0)
	rng.EXPECT().next(-140.0, 140.0).Return(0.0)
	rng.EXPECT().next(-1471.0, 1471.0).Return(0.0)

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

	sut.stop()
}

func TestThatItMovesBetweenCandidates(t *testing.T) {
	objective := NewFakeObjective()
	candidateInput := make(chan *candidate)
	candidateOutput := make(chan *candidate)
	resultOutput := make(chan *candidate)
	random := rand.New(rand.NewSource(6))
	rng := &randwrapper{random}

	sut := newParticle(objective, candidateInput, candidateOutput, resultOutput, rng)
	sut.start()

	localBest := <-objective.Parameters
	objective.Result <- 1
	<-candidateOutput

	globalBest := candidate{make([]float64, 4), 0}
	globalBest.parameters[0] = localBest[0]
	globalBest.parameters[1] = localBest[1] + 1
	globalBest.parameters[2] = localBest[2] + 2
	globalBest.parameters[3] = localBest[3] + 3

	candidateInput <- &globalBest

	candidates := make([][]float64, 10000)

	// swing in and check bounds
	for i := 0; i < 100000; i++ {
		p := <-objective.Parameters
		assert.LessOrEqual(t, p[0], 10.0)
		assert.GreaterOrEqual(t, p[0], 0.0)
		assert.LessOrEqual(t, p[1], 12.0)
		assert.GreaterOrEqual(t, p[1], -1.0)
		// keep result
		objective.Result <- 2
	}

	// sample
	for i := 0; i < 10000; i++ {
		candidates[i] = <-objective.Parameters
		// keep result
		objective.Result <- 2
	}
	avg := calcAverage(candidates)
	std := calcStandardDev(avg, candidates)
	assert.InDelta(t, localBest[0], avg[0], 0.1)
	assert.InDelta(t, localBest[1]+0.5, avg[1], 0.3)
	assert.InDelta(t, localBest[2]+1, avg[2], 0.3)
	assert.InDelta(t, localBest[3]+1.5, avg[3], 1)

	assert.Less(t, std[0], 0.0001)
	assert.Less(t, std[1], 3.0)
	assert.Less(t, std[2], 10.0)
	assert.Less(t, std[3], 20.0)

	sut.stop()
}

func calcStandardDev(avg []float64, candidates [][]float64) []float64 {
	count := len(candidates)
	dim := len(avg)
	result := make([]float64, dim)
	for d := 0; d < dim; d++ {
		for i := 0; i < count; i++ {
			result[d] += math.Pow(candidates[i][d]-avg[d], 2)
		}
		result[d] /= float64(count)
		result[d] = math.Sqrt(result[d])
	}
	return result
}

func calcAverage(candidates [][]float64) []float64 {
	count := len(candidates)
	dim := len(candidates[0])
	result := make([]float64, dim)
	for d := 0; d < dim; d++ {
		for i := 0; i < count; i++ {
			result[d] += candidates[i][d]
		}
		result[d] /= float64(count)
	}
	return result
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
	return float64(10 + 2*math.Pow(float64(dimension), 6))
}

func (f *fakeObjective) Evaluate(parameter []float64) float64 {
	f.Parameters <- parameter
	return <-f.Result
}

var _ Objective = &fakeObjective{}
