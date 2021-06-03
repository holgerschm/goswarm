package goswarm

import "time"

type swarm struct {
	objective   Objective
	topology    topology
	terminators []terminator
}

func (s *swarm) optimize() candidate {
	inputChannels := make([]chan *candidate, s.topology.particleCount())
	globalOutput := make(chan *candidate, 8)
	for i := 0; i < s.topology.particleCount(); i++ {
		inputChannels[i] = make(chan *candidate, 10)
	}
	particles := make([]*particle, s.topology.particleCount())
	for i := 0; i < s.topology.particleCount(); i++ {
		var outputs []chan<- *candidate
		outputs = append(outputs, globalOutput)
		outputMap := s.topology.getOutputs(i)
		for j := 0; j < len(outputMap); j++ {
			outputs = append(outputs, inputChannels[outputMap[j]])
		}
		particles[i] = newParticle(s.objective, inputChannels[i], &nonBlockingMultiplexer{outputs: outputs}, NewSystemRandom(), time.Second)
		particles[i].start()
	}
	for {
		cand := <-globalOutput
		for _, term := range s.terminators {
			if term.shouldTerminate(cand) {
				return *cand
				for i := 0; i < len(particles); i++ {
					particles[i].stop()
				}
			}
		}
	}
}

func newSwarm(obj Objective, top topology, terminators []terminator) *swarm {
	return &swarm{objective: obj, topology: top, terminators: terminators}
}
