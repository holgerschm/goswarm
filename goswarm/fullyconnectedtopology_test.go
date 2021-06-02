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
		assert.Len(t, sut.getOutputs(i), 4)
		for j := 0; j < 5; j++ {
			if j == i {
				assert.NotContains(t, sut.getOutputs(i), j)
			} else {
				assert.Contains(t, sut.getOutputs(i), j)
			}
		}
	}
}
