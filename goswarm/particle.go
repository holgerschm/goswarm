package goswarm

import (
	"math"
	"time"
)

type particle struct {
	objective      Objective
	candidateInput <-chan *Candidate
	output         multiplexer
	rng            random
	position       []float64
	velocity       []float64
	lowerBound     []float64
	upperBound     []float64
	globalBest     *Candidate
	best           *Candidate
	dim            int
	stopping       bool
	stopped        chan bool
	iteration      int64
	updateInterval time.Duration
}

func newParticle(objective Objective, candidateInput <-chan *Candidate, output multiplexer, rng random, updateInterval time.Duration) *particle {
	return &particle{objective,
		candidateInput,
		output,
		rng,
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		nil,
		nil,
		objective.Dimensions(),
		false,
		make(chan bool),
		1,
		updateInterval,
	}
}

func (p *particle) start() {
	go p.run()
}

func (p *particle) run() {
	defer close(p.stopped)
	for i := 0; i < p.objective.Dimensions(); i++ {
		lower := p.objective.GetLowerBound(i)
		upper := p.objective.GetUpperBound(i)
		p.lowerBound[i] = lower
		p.upperBound[i] = upper
		p.position[i] = p.rng.next(lower, upper)
		p.velocity[i] = p.rng.next(lower-upper, upper-lower)
	}
	p.best = p.evaluateCurrent()
	p.globalBest = p.best

	p.output.send(p.best)

	globalIn := <-p.candidateInput
	p.updateGlobalBest(globalIn)

	ticker := time.NewTicker(p.updateInterval)
	defer ticker.Stop()

	candidate := p.best

	for !p.stopping {
		select {
		case globalCandidate := <-p.candidateInput:
			p.updateGlobalBest(globalCandidate)
		case <-ticker.C:
			p.output.send(candidate)
		default:
			p.updateVelocity()
			p.updatePosition()
			p.clampPosition()
			candidate = p.evaluateCurrent()
			p.updateBest(candidate)
			p.updateGlobalBest(candidate)
		}
	}
}

func (p *particle) updateBest(candidate *Candidate) {
	if candidate.Value < p.best.Value {
		p.best = candidate
		p.output.send(candidate)
	}
}

func (p *particle) updateGlobalBest(candidate *Candidate) {
	if candidate.Value < p.globalBest.Value {
		p.globalBest = candidate
	}
}

func (p *particle) evaluateCurrent() *Candidate {
	current := Candidate{}
	current.Parameters = make([]float64, p.dim)
	copy(current.Parameters, p.position)
	current.Value = p.objective.Evaluate(current.Parameters)
	current.Iteration = p.iteration
	p.iteration++
	return &current
}

func (p *particle) updateVelocity() {
	for i := 0; i < p.dim; i++ {
		p.velocity[i] *= 0.72985
		p.velocity[i] += p.rng.next(0, 2.05)*(p.best.Parameters[i]-p.position[i]) +
			p.rng.next(0, 2.05)*(p.globalBest.Parameters[i]-p.position[i])
	}
}

func (p *particle) clampVelocity() {
	for i := 0; i < p.dim; i++ {
		diff := p.upperBound[i] - p.lowerBound[i]
		p.velocity[i] = math.Max(p.velocity[i], -diff)
		p.velocity[i] = math.Min(p.velocity[i], diff)
	}
}

func (p *particle) updatePosition() {
	for i := 0; i < p.dim; i++ {
		p.position[i] += p.velocity[i]
	}
}

func (p *particle) clampPosition() {
	for i := 0; i < p.dim; i++ {
		p.position[i] = math.Max(p.position[i], p.lowerBound[i])
		p.position[i] = math.Min(p.position[i], p.upperBound[i])
	}
}

func (p *particle) stop() {
	p.stopping = true
}

func (p *particle) waitForFinish() {
	<-p.stopped
}
