/*
	Stochastic Gradient Descent

	https://en.wikipedia.org/wiki/Stochastic_gradient_descent

	Minimize a function of the form:
		Sum_i{F_i(x)}, i := 0...terms
*/
package sgd

import (
	"fmt"
	"github.com/deboshire/exp/math/opt"
	"github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/tracer"
	"math"
	"math/rand"
)

type ObjectiveFunc func(x vector.F64) (value float64, gradient vector.F64, ok bool)

type Minimizer struct {
	F      ObjectiveFunc
	X0     vector.F64
	Tracer tracer.Tracer
	State  *State
}

type State struct {
	opt.State

	// Tracer to use
	Tracer tracer.Tracer
}

// Termination criterion generates a double error. The error is compared to eps
// passed to Minimize function and as soon as it is less than eps, optimization
// process is terminated.
// todo(mike): this type name is possibly too long.
type TermCrit interface {
	opt.TermCrit
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

	x0 := minimizer.X0
	if x0 == nil {
		panic(fmt.Errorf("Initial value not set: %q", minimizer))
	}

	minimizer.State = &State{State: *opt.NewState(x0), Tracer: t}
}

func (minimizer *Minimizer) Minimize(eps float64, termCrit TermCrit) (coords vector.F64) {
	minimizer.initIfNeeded()

	s := minimizer.State
	t := s.Tracer.Algorithm("sgd")

	t.TraceFloat64("eps", eps)

	return opt.ConvergeWithState(func(x_param vector.F64, x1 *vector.F64) {
		// One convergence step should is one interation over entire dataset.
		x1.CopyFrom(x_param)

		x := *x1
		f := minimizer.F
		for i := 0; ; i++ {
			// todo(mike): there's some theory about choosing alpha.
			// http://leon.bottou.org/slides/largescale/lstut.pdf
			alpha := .1 / (1 + math.Sqrt(float64(s.TotalIter)))
			//alpha := 1.0 / float64(s.TotalIter)

			y, grad, ok := f(x)
			if !ok {
				// end of data
				return
			}

			grad.Mul(-alpha)
			x.Add(grad)

			t.TraceFloat64("alpha", alpha)
			t.TraceF64("x", x)
			t.TraceF64("grad", grad)
			t.TraceFloat64("y", y)
		}
	}, &s.State, termCrit, eps)
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
		if i >= length {
			i = i % length
			return 0, nil, false
		}
		idx := perm[i]
		i = i + 1

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

	return &Minimizer{F: f, X0: vector.Zeroes(dim)}
}
