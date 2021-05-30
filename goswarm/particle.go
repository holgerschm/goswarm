package goswarm

type particle struct {
	objective       Objective
	candidateInput  <-chan *candidate
	candidateOutput chan<- *candidate
	resultOutput    chan<- *candidate
	rng             random
}

func (p *particle) start() {
	go p.run()
}

func (p *particle) run() {
	current := candidate{}
	dim := p.objective.Dimensions()
	current.parameters = make([]float64, dim)
	for i := 0; i < p.objective.Dimensions(); i++ {
		lower := p.objective.GetLowerBound(i)
		upper := p.objective.GetUpperBound(i)
		current.parameters[i] = p.rng.next(lower, upper)
	}
	p.objective.Evaluate(current.parameters)
}
