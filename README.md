# goswarm

![build & test](https://github.com/holgerschm/goswarm/actions/workflows/go.yml/badge.svg)

This is an implementation of the particle swarm optimization (PSO) method for golang. Each particle of the swarm runs concurrently on its own goroutine.
It can be used to find the global minimum of an n-dimensional function.
The PSO method is suitable for functions with lots of local minima that cannot be treated effectively by gradiend based methods.

## Install
`> go get github.com/holgerschm/goswarm`

## Examples

Define your function with bounds for all parameter dimensions:
````
type testFunction struct {
}

func (t testFunction) Dimensions() int {
    return 2
}

func (t testFunction) Evaluate(parameter []float64) float64 {
    x := parameter[0]
    y := parameter[1]
    return x * x + y * y
}

func (t testFunction) GetLowerBound(dimension int) float64 {
    return -10
}

func (t testFunction) GetUpperBound(dimension int) float64 {
    return 10
}
````
Run the swarm:
````
func main() {
	swarm := goswarm.NewSwarmBuilder(&testFunction{}).Build()
	result := swarm.Minimize()

	fmt.Println("Result after", result.Iteration, "iterations:")
	fmt.Println("Parameters:", result.Parameters)
	fmt.Println("Value:", result.Value)
}
````

Or run with custom configuration:
````
func main() {
	swarm := goswarm.NewSwarmBuilder(&testFunction{}).
		TerminateAfterIterations(100).
		TerminateWhenBelowLimit(0.001).
		WithParticleCount(45).
		WithRingTopology().
		LogTo(&goswarm.ConsoleLogger{}).
		Build()
	result := swarm.Minimize()

	fmt.Println("Result after", result.Iteration, "iterations:")
	fmt.Println("Parameters:", result.Parameters)
	fmt.Println("Value:", result.Value)
}
````
## Options

Topology:

- Fully connected: Every particle in the swarm communicates with every other particle
- Ring: Every particle in the swarm communicates with two other particles.

Terminate:
- After iterations: The iterations for a single particle after that the optimization should stop. This limit is not strictly enforced but means that a single particle will do at least this number of iterations.
- On limit: When a result below the given limit is found the optimization will stop.
- If both termination criteria are given the optimization will stop when the first is met.

Particle count:
- The number of particles the algorithm uses. This is equal to the used threads.

Logger:
- A custom logger can be specified that will be called whenever a lower value is found. The provided `&ConsoleLogger{}` will print everything to the console, which can be a lot. You can implement your own logger to throttle this.

## Thread safety

The given objective function has to be thread safe as it will be called from multiple threads.



