package stat

import (
	"math"
)

// Window holds last N values and can calculate various statistics on them.
type Window struct {
	N      int
	values []float64
}

func (w *Window) Add(v float64) {
	w.values = append(w.values, v)
	if len(w.values) > w.N {
		w.values = w.values[1:]
	}
}

func (w *Window) Full() bool {
	return len(w.values) == w.N
}

func (w *Window) Max() float64 {
	result := math.SmallestNonzeroFloat64
	for _, v := range w.values {
		if v > result {
			result = v
		}
	}

	return result
}

func (w *Window) Mean() float64 {
	result := 0.0
	for _, v := range w.values {
		result += v
	}

	return result / float64(len(w.values))
}
