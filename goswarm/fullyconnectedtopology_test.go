package goswarm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatItConnectsAllToAll(t *testing.T) {
	var sut topology
	sut = newFullyConnectedTopology(5)

	assert.Equal(t, 5, sut.particleCount())
	for i := 0; i < 5; i++ {
		assert.Len(t, sut.getOutputs(i), 5)
		for j := 0; j < 5; j++ {
			assert.Equal(t, j, sut.getOutputs(i)[j])
		}
	}
}
