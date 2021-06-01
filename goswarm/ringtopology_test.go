package goswarm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatItConnectsNeighbours(t *testing.T) {
	var sut topology
	sut = newRingTopology(5)

	assert.Equal(t, 5, sut.particleCount())
	for i := 0; i < 5; i++ {
		assert.Len(t, sut.getOutputs(i), 2)
		assert.Equal(t, (i+4)%5, sut.getOutputs(i)[0])
		assert.Equal(t, (i+6)%5, sut.getOutputs(i)[1])
	}
}
