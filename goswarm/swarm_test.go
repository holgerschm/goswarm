package goswarm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatItFindsAMinimumIn1D(t *testing.T) {
	var obj Objective
	obj = &parabola{}
	t1 := &iterationsTerminator{minIterations: 100000}
	t2 := &targetValueTerminator{-100} // unreachable

	sut := newSwarm(obj, newFullyConnectedTopology(20), []terminator{t1, t2})

	var result candidate
	result = sut.optimize()

	assert.InDelta(t, 2, result.parameters[0], 0.001)
	assert.InDelta(t, -1, result.value, 0.001)
}

func TestThatItBreaksOnMinimumValue(t *testing.T) {
	var obj Objective
	obj = &parabola{}
	t1 := &iterationsTerminator{minIterations: 1000000000}
	t2 := &targetValueTerminator{10} // unreachable

	sut := newSwarm(obj, newFullyConnectedTopology(20), []terminator{t1, t2})

	var result candidate
	result = sut.optimize()

	assert.Less(t, result.value, 10.0)
	assert.Less(t, result.iteration, int64(1000))
}
