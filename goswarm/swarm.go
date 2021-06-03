package goswarm

import "time"

type Swarm interface {
	Minimize() Candidate
}

type swarm struct {
	objective   Objective
	topology    topology
	terminators []terminator
	logger      Logger
}

func (s *swarm) Minimize() Candidate {
	inputChannels := make([]chan *Candidate, s.topology.particleCount())
	globalOutput := make(chan *Candidate, 8)
	for i := 0; i < s.topology.particleCount(); i++ {
		inputChannels[i] = make(chan *Candidate, 10)
	}
	particles := make([]*particle, s.topology.particleCount())
	for i := 0; i < s.topology.particleCount(); i++ {
		var outputs []chan<- *Candidate
		outputs = append(outputs, globalOutput)
		outputMap := s.topology.getOutputs(i)
		for j := 0; j < len(outputMap); j++ {
			outputs = append(outputs, inputChannels[outputMap[j]])
		}
		particles[i] = newParticle(s.objective, inputChannels[i], &nonBlockingMultiplexer{outputs: outputs}, NewSystemRandom(), 250*time.Millisecond)
		particles[i].start()
	}
	best := <-globalOutput
	for {
		cand := <-globalOutput
		if cand.Value < best.Value {
			best = cand
			s.logger.Log(*best)
		}
		for _, term := range s.terminators {
			if term.shouldTerminate(cand) {
				for i := 0; i < len(particles); i++ {
					particles[i].stop()
				}
				for i := 0; i < len(particles); i++ {
					particles[i].waitForFinish()
				}
				return *best
			}
		}
	}
}

func newSwarm(obj Objective, top topology, terminators []terminator, logger Logger) *swarm {
	return &swarm{objective: obj, topology: top, terminators: terminators, logger: logger}
}
