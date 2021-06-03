package goswarm

type Builder struct {
	particleCount  int
	ringTopology   bool
	iterationLimit int64
	limit          *float64
	logger         Logger
	objective      Objective
}

func (b *Builder) WithFullyConnectedTopology() *Builder {
	b.ringTopology = false
	return b
}

func (b *Builder) WithRingTopology() *Builder {
	b.ringTopology = true
	return b
}

func (b *Builder) WithParticleCount(count int) *Builder {
	b.particleCount = count
	return b
}

func (b *Builder) TerminateAfterIterations(iterations int64) *Builder {
	b.iterationLimit = iterations
	return b
}

func (b *Builder) TerminateWhenBelowLimit(limit float64) *Builder {
	b.limit = &limit
	return b
}

func (b *Builder) Build() Swarm {
	var topology topology
	if b.ringTopology {
		topology = newRingTopology(b.particleCount)
	} else {
		topology = newFullyConnectedTopology(b.particleCount)
	}
	var terminators []terminator
	if b.iterationLimit > 0 {
		terminators = append(terminators, &iterationsTerminator{b.iterationLimit})
	}
	if b.limit != nil {
		terminators = append(terminators, &targetValueTerminator{*b.limit})
	}

	if len(terminators) == 0 {
		terminators = append(terminators, &iterationsTerminator{10000})
	}

	return newSwarm(b.objective, topology, terminators, b.logger)
}

func (b *Builder) LogTo(logger Logger) *Builder {
	b.logger = logger
	return b
}

func NewSwarmBuilder(obj Objective) *Builder {
	return &Builder{objective: obj, particleCount: 20, ringTopology: true, iterationLimit: 0, logger: &nilLogger{}}
}
