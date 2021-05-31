package goswarm

import "math"

type particle struct {
	objective       Objective
	candidateInput  <-chan *candidate
	candidateOutput chan<- *candidate
	resultOutput    chan<- *candidate
	rng             random
	position        []float64
	velocity        []float64
	lowerBound      []float64
	upperBound      []float64
	globalBest      *candidate
	best            *candidate
	dim             int
	stopped         bool
}

func newParticle(objective Objective,
	candidateInput <-chan *candidate,
	candidateOutput chan<- *candidate,
	resultOutput chan<- *candidate,
	rng random) *particle {
	return &particle{objective,
		candidateInput,
		candidateOutput,
		resultOutput,
		rng,
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		make([]float64, objective.Dimensions()),
		nil,
		nil,
		objective.Dimensions(),
		false}
}

func (p *particle) start() {
	go p.run()
}

func (p *particle) run() {
	for i := 0; i < p.objective.Dimensions(); i++ {
		lower := p.objective.GetLowerBound(i)
		upper := p.objective.GetUpperBound(i)
		p.lowerBound[i] = lower
		p.upperBound[i] = upper
		p.position[i] = p.rng.next(lower, upper)
		p.velocity[i] = p.rng.next(lower-upper, upper-lower)
	}
	p.best = p.evaluateCurrent()

	p.candidateOutput <- p.best

	p.globalBest = <-p.candidateInput

	for !p.stopped {
		p.updateVelocity()
		p.updatePosition()
		p.clampPosition()
		p.evaluateCurrent()
	}
}

func (p *particle) evaluateCurrent() *candidate {
	current := candidate{}
	current.parameters = make([]float64, p.dim)
	copy(current.parameters, p.position)
	current.value = p.objective.Evaluate(current.parameters)
	return &current
}

func (p *particle) updateVelocity() {
	for i := 0; i < p.dim; i++ {
		p.velocity[i] *= 0.8
		p.velocity[i] += p.rng.next(0, 2)*(p.best.parameters[i]-p.position[i]) +
			p.rng.next(0, 2)*(p.globalBest.parameters[i]-p.position[i])
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
	p.stopped = true
}
