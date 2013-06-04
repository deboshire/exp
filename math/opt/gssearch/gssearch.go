// Golden Section Search
// https://en.wikipedia.org/wiki/Golden_section_search
package gssearh

import (
	"math"
)

var (
	phi    = (1 + math.Sqrt(5)) / 2
	resphi = 2 - phi
)

type State struct {
	// iteration
	iter int

	// coordinates
	A, B, C, X float64

	// function applied to them
	FB, FX float64
}

type TermCrit interface {
	ShouldTerminate(s *State) float64
}

type AbsoluteErrorTermCrit struct{}

func (c *AbsoluteErrorTermCrit) ShouldTerminate(s *State) float64 {
	return math.Abs(s.C-s.A) / (math.Abs(s.B) + math.Abs(s.X))
}

func Minimize(minX float64, maxX float64, f func(float64) float64, termCrit TermCrit, eps float64) (res float64) {
	b := minX + resphi*(maxX-minX)
	return minimize(State{A: minX, B: b, C: maxX, FB: f(b)}, f, termCrit, eps)
}

func minimize(state State, f func(float64) float64, termCrit TermCrit, eps float64) (res float64) {
	for iter := 0; ; iter++ {
		state.iter = iter
		a := state.A
		b := state.B
		c := state.C

		if c-b > b-a {
			state.X = b + resphi*(c-b)
		} else {
			state.X = b - resphi*(b-a)
		}

		x := state.X
		state.FX = f(x)

		if termCrit.ShouldTerminate(&state) < eps {
			return (a + c) / 2
		}

		if state.FX < state.FB {
			if c-b > b-a {
				state.A = b
				state.B = x
				state.FB = state.FX
			} else {
				state.B = x
				state.C = b
				state.FB = state.FX
			}
		} else {
			if c-b > b-a {
				state.C = x
			} else {
				state.A = x
			}
		}
	}
}
