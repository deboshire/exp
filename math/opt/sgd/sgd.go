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
)

////////////////////////////////////////////////////////////////////////////////

type PushFunction func(value float64, gradient vector.F64)
type IterationFn func(pushFn PushFunction)

type PushMinimizer struct {
	X0     vector.F64
	State  *State
	Tracer tracer.Tracer
}

// TODO: dedupe with iteration minimizer
func (pushMinimizer *PushMinimizer) Minimize(eps float64, termCrit TermCrit, iterationFn IterationFn) (coords vector.F64) {
	pushMinimizer.initIfNeeded()

	s := pushMinimizer.State
	t := s.Tracer.Algorithm("sgd")

	t.TraceFloat64("eps", eps)

	// TODO: use push converge
	return opt.ConvergeWithState(func(x_param vector.F64, x1 *vector.F64) {
		x1.CopyFrom(x_param)
		x := *x1

		pushFn := func(value float64, grad vector.F64) {
			alpha := .1 / (1 + math.Sqrt(float64(s.TotalIter)))
			grad.Mul(-alpha)
			x.Add(grad)
		}
		iterationFn(pushFn)
	}, &s.State, termCrit, eps)
}

// TODO: dedupe
func (pushMinimizer *PushMinimizer) initIfNeeded() {
	if pushMinimizer.State != nil {
		return
	}

	t := pushMinimizer.Tracer
	if t == nil {
		t = tracer.DefaultTracer()
	}

	x0 := pushMinimizer.X0
	if x0 == nil {
		panic(fmt.Errorf("Initial value not set: %q", pushMinimizer))
	}

	pushMinimizer.State = &State{State: *opt.NewState(x0), Tracer: t}
}

//////////////////////////////////////////////////////////////////////////////////

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
