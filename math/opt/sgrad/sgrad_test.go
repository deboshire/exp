package sgrad

import (
	"github.com/deboshire/exp/math/vector"
	"math"
	"math/rand"
	"testing"
)

func TestLeastSquaresPrecise(t *testing.T) {
	minimizer := LeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})

	term := AbsDistanceCrit{NumItersToAvg: 10}
	v, coords := minimizer.Minimize(1e-10, &term)
	t.Log("Value: ", v, "Coords: ", coords)

	if math.Abs(coords[0]) > 1e-2 {
		t.Error("coords[0] != 0: %s", coords[0])
	}
	if math.Abs(coords[1]-1) > 1e-2 {
		t.Error("coords[1] != 1: %s", coords[0])
	}
}

func TestLeastSquares_NumIterations(t *testing.T) {
	// Demo problem from https://en.wikipedia.org/wiki/Linear_least_squares_(mathematics)
	minimizer := LeastSquares([]vector.F64{
		vector.F64{1, 6},
		vector.F64{2, 5},
		vector.F64{3, 7},
		vector.F64{4, 10},
	})

	term := NumIterationsCrit{NumIterations: 100000}
	v, coords := minimizer.Minimize(1e-8, &term)
	t.Log("Value: ", v, "Coords: ", coords)
	if math.Abs(coords[0]-3.5) > 1e-1 {
		t.Error("coords[0] != 3.5: %s", coords[0])
	}
	if math.Abs(coords[1]-1.4) > 1e-1 {
		t.Error("coords[1] != 1.4: %s", coords[0])
	}
}

func BenchmarkLeastSquare(b *testing.B) {
	minimizer := LeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})

	term := NumIterationsCrit{NumIterations: b.N}
	minimizer.Minimize(1e-8, &term)
}

func init() {
	rand.Seed(1)
}
