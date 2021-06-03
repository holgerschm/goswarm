package goswarm

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestThatItStartsAtARandomPosition(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objective := NewFakeObjective()
	candidateInput := make(chan *Candidate)

	outputChannel := make(chan *Candidate)
	outputMP := blockingMultiplexer{[]chan<- *Candidate{outputChannel}}

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

	sut := newParticle(objective, candidateInput, outputMP, rng, time.Hour)
	sut.start()

	parameters := <-objective.Parameters
	objective.Result <- 7.7

	assert.Equal(t, 1.1, parameters[0])
	assert.Equal(t, 2.2, parameters[1])
	assert.Equal(t, 3.3, parameters[2])
	assert.Equal(t, 4.4, parameters[3])

	// output first Value as best
	output := <-outputChannel

	assert.Equal(t, 1.1, output.Parameters[0])
	assert.Equal(t, 2.2, output.Parameters[1])
	assert.Equal(t, 3.3, output.Parameters[2])
	assert.Equal(t, 4.4, output.Parameters[3])
	assert.Equal(t, 7.7, output.Value)
	assert.Equal(t, int64(1), output.Iteration)

	sut.stop()
	candidateInput <- &Candidate{}
	sut.waitForFinish()
}

func TestThatItMovesBetweenCandidates(t *testing.T) {
	defer goleak.VerifyNone(t)
	objective := NewFakeObjective()
	candidateInput := make(chan *Candidate, 1)
	outputChannel := make(chan *Candidate)
	outputMP := blockingMultiplexer{[]chan<- *Candidate{outputChannel}}
	random := rand.New(rand.NewSource(6))
	rng := &randwrapper{random}

	sut := newParticle(objective, candidateInput, outputMP, rng, 1000*time.Second)
	sut.start()

	localBest := <-objective.Parameters
	objective.Result <- 1
	<-outputChannel

	globalBest := &Candidate{make([]float64, 4), 0, 0}
	globalBest.Parameters[0] = localBest[0]
	globalBest.Parameters[1] = localBest[1] + 1
	globalBest.Parameters[2] = localBest[2] + 2
	globalBest.Parameters[3] = localBest[3] + 3

	candidateInput <- globalBest

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

	// change global optimum
	globalBest = &Candidate{make([]float64, 4), -1, 0}
	globalBest.Parameters[0] = localBest[0]
	globalBest.Parameters[1] = localBest[1] - 1
	globalBest.Parameters[2] = localBest[2] - 2
	globalBest.Parameters[3] = localBest[3] - 3

	candidateInput <- globalBest
	<-objective.Parameters
	objective.Result <- 2

	// should be ignored
	globalBest = &Candidate{make([]float64, 4), -0.9, 0}
	globalBest.Parameters[0] = localBest[0]
	globalBest.Parameters[1] = localBest[1] + 1
	globalBest.Parameters[2] = localBest[2] + 2
	globalBest.Parameters[3] = localBest[3] + 30
	candidateInput <- globalBest

	// swing in and check bounds
	for i := 0; i < 100000; i++ {
		<-objective.Parameters
		// keep result
		objective.Result <- 2
	}

	// sample
	for i := 0; i < 10000; i++ {
		candidates[i] = <-objective.Parameters
		// keep result
		objective.Result <- 2
	}

	avg = calcAverage(candidates)
	std = calcStandardDev(avg, candidates)
	assert.InDelta(t, localBest[0], avg[0], 0.1)
	assert.InDelta(t, localBest[1]-0.5, avg[1], 0.3)
	assert.InDelta(t, localBest[2]-1, avg[2], 0.3)
	assert.InDelta(t, localBest[3]-1.5, avg[3], 1)

	assert.Less(t, std[0], 0.0001)
	assert.Less(t, std[1], 3.0)
	assert.Less(t, std[2], 10.0)
	assert.Less(t, std[3], 30.0)

	// change local optimum == global optimum
	newOptimum := <-objective.Parameters
	objective.Result <- -2

	out := <-outputChannel
	assert.Equal(t, out.Value, -2.0)
	assert.Equal(t, out.Parameters, newOptimum)
	assert.Equal(t, out.Iteration, int64(220003))

	// swing in and check bounds
	for i := 0; i < 100000; i++ {
		<-objective.Parameters
		// keep result
		objective.Result <- 2
	}

	// sample
	for i := 0; i < 10000; i++ {
		candidates[i] = <-objective.Parameters
		// keep result
		objective.Result <- 2
	}

	avg = calcAverage(candidates)
	std = calcStandardDev(avg, candidates)
	assert.InDelta(t, newOptimum[0], avg[0], 0.1)
	assert.InDelta(t, newOptimum[1], avg[1], 0.1)
	assert.InDelta(t, newOptimum[2], avg[2], 0.1)
	assert.InDelta(t, newOptimum[3], avg[3], 0.1)

	assert.Less(t, std[0], 0.0001)
	assert.Less(t, std[1], 0.0001)
	assert.Less(t, std[2], 0.0001)
	assert.Less(t, std[3], 0.0001)

	sut.stop()
	for {
		select {
		case <-objective.Parameters:
			objective.Result <- 0
		default:
			sut.waitForFinish()
			return
		}
	}
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
	return 10.0 + 2.0*math.Pow(float64(dimension), 6)
}

func (f *fakeObjective) Evaluate(parameter []float64) float64 {
	f.Parameters <- parameter
	return <-f.Result
}

var _ Objective = &fakeObjective{}

type blockingMultiplexer struct {
	outputs []chan<- *Candidate
}

func (b blockingMultiplexer) send(best *Candidate) {
	for i := 0; i < len(b.outputs); i++ {
		b.outputs[i] <- best
	}
}

var _ multiplexer = &blockingMultiplexer{}
