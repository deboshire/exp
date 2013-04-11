/*
	Stochastic Gradient Descent

	https://en.wikipedia.org/wiki/Stochastic_gradient_descent

	Minimize a function of the form:
		Sum_i{F_i(x)}, i := 0...terms
*/
package sgrad

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/tracer"
	"math"
	"math/rand"
)

type ObjectiveFunc func(x vector.F64) (value float64, gradient vector.F64, ok bool)

type Minimizer struct {
	F       ObjectiveFunc
	Initial vector.F64
	Tracer  tracer.Tracer
	State   *State
}

type State struct {
	// Tracer to use
	Tracer tracer.Tracer
	// Epoch of minimizer, i.e. how many times Minimize was called
	Epoch int
	// Number of iterations performed in current epoch
	Iter int
	// Number of iterations performed in previous epochs
	EpochIter []int
	// Total number of iterations performed.
	TotalIter int

	// Value of the function term computed on the last iteration. Quite meaningless
	Value float64

	// Current X vector
	X vector.F64

	// Last change of X vector
	Dx vector.F64
}

// Termination criterion generates a double error. The error is compared to eps
// passed to Minimize function and as soon as it is less than eps, optimization
// process is terminated.
// todo(mike): this type name is possibly too long.
type TermCrit interface {
	ShouldTerminate(s *State) float64
}

/*
// Termination criterion that keeps track of relative mean improvement of the
// function value
type RelativeMeanImprovementCrit struct {
	NumItersToAvg int
	prevVals      []float64
}

func (c *RelativeMeanImprovementCrit) ShouldTerminate(s *State) float64 {
	iters := c.NumItersToAvg
	if iters < 2 {
		iters = 5
	}
	c.prevVals = append(c.prevVals, s.Value)

	if len(c.prevVals) < iters {
		// not enough values yet.
		return math.MaxFloat64
	}

	if len(c.prevVals) > iters {
		c.prevVals = c.prevVals[1:]
	}

	prevVal := c.prevVals[0]
	val := s.Value
	avgImprovement := (prevVal - val) / float64(iters)
	relAvgImpr := math.Abs(avgImprovement / val)
	s.Tracer.TraceFloat64("avgImprovement", avgImprovement)
	s.Tracer.TraceFloat64("relAvgImpr", relAvgImpr)
	return relAvgImpr
}

*/

// DONOTSUBMIT - this is a horrible criterion. Implement another one.
type AbsDistanceCrit struct {
	NumItersToAvg int
	prevVals      []float64
}

func (c *AbsDistanceCrit) ShouldTerminate(s *State) float64 {
	iters := c.NumItersToAvg
	if iters < 2 {
		iters = 5
	}
	val := s.Dx.Length()
	c.prevVals = append(c.prevVals, val)

	if len(c.prevVals) < iters {
		// not enough values yet.
		return math.MaxFloat64
	}

	if len(c.prevVals) > iters {
		c.prevVals = c.prevVals[1:]
	}

	prevVal := c.prevVals[0]
	avgImprovement := (prevVal - val) / float64(iters)
	s.Tracer.TraceFloat64("avgImprovement", avgImprovement)
	return math.Abs(avgImprovement)
}

// Termination criterion that terminates after given number of iterations.
type NumIterationsCrit struct {
	NumIterations int
}

func (c *NumIterationsCrit) ShouldTerminate(s *State) float64 {
	if s.Iter >= c.NumIterations - 1 {
		return 0
	}
	return math.MaxFloat64

}

// sets all fields to default if this is the first call
func (minimizer *Minimizer) initIfNeeded() {
	if minimizer.State != nil {
		return
	}

	t := minimizer.Tracer
	if t == nil {
		t = tracer.DefaultTracer()
	}

	x := minimizer.Initial
	if x == nil {
		panic(fmt.Errorf("Initial value not set: %q", minimizer))
	}

	// First run.
	minimizer.State = &State{X: x, Tracer: t}
}

func (minimizer *Minimizer) Minimize(eps float64, term TermCrit) (value float64, coords vector.F64) {
	minimizer.initIfNeeded()

	s := minimizer.State
	x := s.X
	t := s.Tracer
	f := minimizer.F

	s.Epoch++

	for i := 0; ; i++ {
		s.Iter = i
		s.TotalIter++

		// todo(mike): there's some theory about choosing alpha.
		// http://leon.bottou.org/slides/largescale/lstut.pdf
		alpha := .1 / (1 + math.Sqrt(float64(i)))

		t.TraceFloat64("alpha", alpha)

		t.TraceF64("x", x)

		y, grad, ok := f(x)
		if !ok {
			// end of data
			break
		}
		s.Value = y

		t.TraceF64("grad", grad)
		t.TraceFloat64("y", y)

		grad.Mul(-alpha)
		x.Add(grad)

		s.Dx = grad
		s.X = x

		err := term.ShouldTerminate(s)
		t.TraceFloat64("err", err)
		//fmt.Println(i, ",", alpha, ",", x, ",", grad, ",", y, ",", err)
		if err < eps {
			break
		}
	}

	s.EpochIter = append(s.EpochIter, s.Iter)
	return s.Value, s.X
}

/*
	Objective function for performing least squares optimization.
*/
func LeastSquares(points []vector.F64) *Minimizer {
	dim := len(points[0])
	length := len(points)

	perm := rand.Perm(len(points))
	grad := vector.Zeroes(dim)
	i := 0

	f := func(x vector.F64) (value float64, gradient vector.F64, ok bool) {
		idx := perm[i]
		i = (i + 1) % length

		// The function itself is:
		// (x[0] + x[1]*row[0] + x[2]*row[1] + .... - row)^2
		a := x[0]
		row := points[idx]
		for i := 0; i < dim-1; i++ {
			a += x[i+1] * row[i]
		}
		a -= row[dim-1]

		// The gradient is
		// 2a for i == 0, 2*points[idx][i-1]*a for other idx
		grad[0] = 2 * a
		for i := 1; i < dim; i++ {
			grad[i] = 2 * row[i-1] * a
		}
		return a*a*0.5 + 1, grad, true
	}

	return &Minimizer{F: f, Initial: vector.Zeroes(dim)}
}
