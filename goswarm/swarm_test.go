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

	builder := NewSwarmBuilder(obj).WithFullyConnectedTopology()
	sut := builder.WithParticleCount(20).
		TerminateAfterIterations(100000).
		TerminateWhenBelowLimit(-100). // unreachable
		Build()

	var result Candidate
	result = sut.Minimize()

	assert.InDelta(t, 2, result.Parameters[0], 0.001)
	assert.InDelta(t, -1, result.Value, 0.001)

	time.Sleep(100 * time.Millisecond)
}

func TestThatItBreaksOnMinimumValue(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &parabola{}

	sut := NewSwarmBuilder(obj).
		TerminateWhenBelowLimit(10).
		TerminateAfterIterations(1000000000).
		Build()

	var result Candidate
	result = sut.Minimize()

	assert.Less(t, result.Value, 10.0)
	assert.Less(t, result.Iteration, int64(1000000))

	time.Sleep(100 * time.Millisecond)
}

func TestThatItFindsAGlobalMinimumIn2D(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &levi{}

	sut := NewSwarmBuilder(obj).
		WithRingTopology().
		TerminateAfterIterations(1000000).
		Build()

	var result Candidate
	result = sut.Minimize()

	assert.InDelta(t, 1, result.Parameters[0], 0.001)
	assert.InDelta(t, 1, result.Parameters[1], 0.001)
	assert.InDelta(t, 0, result.Value, 0.001)

	time.Sleep(100 * time.Millisecond)
}

func TestThatItFindsAGlobalMinimumInHigherDimensions(t *testing.T) {
	defer goleak.VerifyNone(t)

	var obj Objective
	obj = &rastrigin{}

	sut := NewSwarmBuilder(obj).
		WithRingTopology().
		WithParticleCount(40).
		TerminateAfterIterations(1000000).
		Build()

	var result Candidate
	result = sut.Minimize()

	for i := 0; i < obj.Dimensions(); i++ {
		assert.InDelta(t, 0, result.Parameters[i], 0.001)
	}
	assert.InDelta(t, 0, result.Value, 0.001)

	time.Sleep(100 * time.Millisecond)
}
