package opt

import (
	"github.com/deboshire/exp/math/stat"
	"github.com/deboshire/exp/math/vector"
	"math"
)

// TermCrit is a termination criterion for iterations.
// Algorithm terminates when ShouldTerminate becomes < epsilon.
type TermCrit interface {
	ShouldTerminate(s *State) float64
}

type State struct {
	Epoch     int
	Iter      int
	TotalIter int
	// previous X value
	X vector.F64
	// new X value
	X1 vector.F64

	// Number of iterations performed in previous epochs
	EpochIter []int
}

// Converge applies the function until it converges.
// f takes the old vectox X and stores result into x1.
func Converge(f func(x vector.F64, x1 *vector.F64), x0 vector.F64, termCrit TermCrit, eps float64) vector.F64 {
	return ConvergeWithState(f, NewState(x0), termCrit, eps)
}

func NewState(x0 vector.F64) *State {
	return &State{X: x0.Copy(), X1: x0.Copy()}
}

func ConvergeWithState(f func(x vector.F64, x1 *vector.F64), state *State, termCrit TermCrit, eps float64) vector.F64 {
	pushFn := PushConvergeWithState(state, termCrit, eps)

	x1 := state.X.Copy()
	for {
		f(state.X, &x1)
		ok := pushFn(x1)
		if !ok {
			return state.X
		}
	}
}

type PushConvergeFn func(x1 vector.F64) bool

func PushConverge(x0 vector.F64, termCrit TermCrit, eps float64) PushConvergeFn {
	return PushConvergeWithState(NewState(x0), termCrit, eps)
}

func PushConvergeWithState(state *State, termCrit TermCrit, eps float64) PushConvergeFn {
	state.Epoch++
	state.Iter = 0

	return func(x1 vector.F64) bool {
		state.X1.CopyFrom(x1)

		// fmt.Println("x:", state.X, "x1: ", state.X1)

		if termCrit.ShouldTerminate(state) <= eps {
			state.EpochIter = append(state.EpochIter, state.Iter+1)
			return false
		}

		temp := state.X
		state.X = state.X1
		state.X1 = temp
		state.Iter++
		state.TotalIter++
		return true
	}
}

// Termination criterion that terminates after given number of iterations.
type NumIterationsCrit struct {
	NumIterations int
}

func (c *NumIterationsCrit) ShouldTerminate(s *State) float64 {
	if s.Iter >= c.NumIterations {
		return 0
	}
	return math.MaxFloat64
}

// Tracks mean of relative coordinate change vector length.
type MaxRelativeChangeCrit struct {
	Iters  int
	window *stat.Window
}

func (c *MaxRelativeChangeCrit) ShouldTerminate(s *State) float64 {
	if c.window == nil {
		if c.Iters < 1 {
			c.Iters = 10
		}

		c.window = &stat.Window{N: c.Iters}
	}

	dx := s.X.Copy()
	dx.Sub(s.X1)
	dxl := dx.Length()
	sxl := s.X.Length()

	val := dxl / sxl

	c.window.Add(val)

	if !c.window.Full() {
		return math.MaxFloat64
	}

	// fmt.Println("dxl:", dxl, "sxl:", sxl, "val:", val)
	return c.window.Max()
}

// Tracks mean of relative coordinate change vector length.
type MaxAbsChangeCrit struct {
	Iters  int
	window *stat.Window
}

func (c *MaxAbsChangeCrit) ShouldTerminate(s *State) float64 {
	if c.window == nil {
		if c.Iters < 1 {
			c.Iters = 10
		}

		c.window = &stat.Window{N: c.Iters}
	}

	dx := s.X.Copy()
	dx.Sub(s.X1)
	dxl := dx.Length()
	c.window.Add(dxl)

	if !c.window.Full() {
		return math.MaxFloat64
	}

	return c.window.Max()
}
