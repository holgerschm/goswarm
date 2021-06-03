// all tests in this file are integration test that are non-deterministic, i.e. they may fail from time to time

package goswarm

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func TestThatItFindsAMinimumIn1D(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &parabola{}
	t1 := &iterationsTerminator{minIterations: 100000}
	t2 := &targetValueTerminator{-100} // unreachable

	sut := newSwarm(obj, newFullyConnectedTopology(20), []terminator{t1, t2})

	var result candidate
	result = sut.optimize()

	assert.InDelta(t, 2, result.parameters[0], 0.001)
	assert.InDelta(t, -1, result.value, 0.001)

	time.Sleep(100 * time.Millisecond)
}

func TestThatItBreaksOnMinimumValue(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &parabola{}
	t1 := &iterationsTerminator{minIterations: 1000000000}
	t2 := &targetValueTerminator{10} // unreachable

	sut := newSwarm(obj, newFullyConnectedTopology(20), []terminator{t1, t2})

	var result candidate
	result = sut.optimize()

	assert.Less(t, result.value, 10.0)
	assert.Less(t, result.iteration, int64(1000000))

	time.Sleep(100 * time.Millisecond)
}

func TestThatItFindsAGlobalMinimumIn2D(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &levi{}
	t1 := &iterationsTerminator{minIterations: 1000000}

	sut := newSwarm(obj, newRingTopology(20), []terminator{t1})

	var result candidate
	result = sut.optimize()

	assert.InDelta(t, 1, result.parameters[0], 0.001)
	assert.InDelta(t, 1, result.parameters[1], 0.001)
	assert.InDelta(t, 0, result.value, 0.001)

	time.Sleep(100 * time.Millisecond)
}

func TestThatItFindsAGlobalMinimumInHigherDimensions(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &rastrigin{}
	t1 := &iterationsTerminator{minIterations: 1000000}

	sut := newSwarm(obj, newRingTopology(40), []terminator{t1})

	var result candidate
	result = sut.optimize()

	for i := 0; i < obj.Dimensions(); i++ {
		assert.InDelta(t, 0, result.parameters[i], 0.001)
	}
	assert.InDelta(t, 0, result.value, 0.001)

	time.Sleep(100 * time.Millisecond)
}
