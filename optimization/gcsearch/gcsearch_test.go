package gcsearh

import (
	"math"
	"testing"
)

func TestGoldenCutSearch(t *testing.T) {
	f := func(x float64) float64 {
		return 5*x*x - 4*x - 3
	}

	x := Minimize(-10, 10, f, &AbsoluteErrorTermCrit{}, 1e-10)

	if math.Abs(x - 0.4) > 1e-8 {
		t.Errorf("x=%f", x)
	}
}
