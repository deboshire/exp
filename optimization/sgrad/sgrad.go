/*
	Stochastic Gradient Descent

	https://en.wikipedia.org/wiki/Stochastic_gradient_descent
*/
package sgrad

import (
	"github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/tracer"
	"math"
	"math/rand"
)

type ObjectiveFunc struct {
	Terms int
	F     func(idx int, x vector.F64, out_gradient vector.F64) float64
}

type State struct {
	Tracer tracer.Tracer
	Pass   int
	Value  float64
}

// Termination criterion generates a double error. The error is compared to epsilon
// passed to Minimize function and as soon as it is less than epsilon, optimization
// process is terminated.
// todo(mike): this type name is possibly too long.
type TerminationCriterion interface {
	ShouldTerminate(s *State) float64
}

// Termination criterion that keeps track of relative mean improvement of the
// function value
type RelativeMeanImprovementCriterion struct {
	NumItersToAvg int
	prevVals      []float64
}

func (c *RelativeMeanImprovementCriterion) ShouldTerminate(s *State) float64 {
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

// Termination criterion that terminates after given number of iterations.
type NumIterationsCriterion struct {
	NumIterations int
}

func (c *NumIterationsCriterion) ShouldTerminate(s *State) float64 {
	if s.Pass >= c.NumIterations-1 {
		return 0
	}
	return math.MaxFloat64

}

/*
	Minimize a function of the form:
		Sum_i{F_i(x)}, i := 0...terms
*/
func Minimize(f ObjectiveFunc, initial vector.F64, eps float64, term TerminationCriterion, t tracer.Tracer) (value float64, coords vector.F64) {
	if t == nil {
		t = tracer.DefaultTracer()
	}

	s := State{Pass: 0, Tracer: t}
	x := initial.Copy()
	grad := initial.Copy()

	for pass := 0; ; pass++ {
		s.Pass = pass
		perm := rand.Perm(f.Terms)
		maxDist := 0.0

		// todo(mike): there's some theory about choosing alpha.
		// http://leon.bottou.org/slides/largescale/lstut.pdf
		alpha := .1 / (1 + math.Sqrt(float64(pass)))
		t.TraceFloat64("alpha", alpha)

		for _, idx := range perm {
			t.TraceInt("idx", idx)

			t.TraceF64("x", x)

			y:= f.F(idx, x, grad)
			t.TraceF64("grad", grad)
			t.TraceFloat64("y", y)

			grad.Mul(-alpha)
			grad.Add(x)

			dist := x.Dist2(grad)
			if dist > maxDist {
				maxDist = dist
			}
			temp := x
			x = grad
			grad = temp
			value = y
		}

		t.TraceFloat64("maxDist", maxDist)

		s.Value = value
		err := term.ShouldTerminate(&s)
		t.TraceFloat64("err", err)
		if err < eps {
			break
		}
	}

	return value, x
}

/*
	Objective function for performing least squares optimization.
*/
func LeastSquares(points []vector.F64) ObjectiveFunc {
	dim := len(points[0])

	f := func(idx int, x vector.F64, gradient vector.F64) (value float64) {
		// The function itself is:
		// (x[0] + x[1]*points[idx][0] + x[2]*points[idx][1] + .... - points[-1])^2
		a := x[0]
		row := points[idx]
		for i := 0; i < dim-1; i++ {
			a += x[i+1] * row[i]
		}
		a -= row[dim-1]

		// The gradient is
		// 2a for i == 0, 2*points[idx][i-1]*a for other idx
		gradient[0] = 2 * a
		for i := 1; i < dim; i++ {
			gradient[i] = 2 * row[i-1] * a
		}

		value = a*a*0.5 + 1
		return
	}

	return ObjectiveFunc{Terms: len(points), F: f}
}
