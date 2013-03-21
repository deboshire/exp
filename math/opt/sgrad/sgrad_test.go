package sgrad

import (
	"github.com/deboshire/exp/math/vector"
	"math"
	"math/rand"
	"testing"
)

func TestLeastSquaresPrecise(t *testing.T) {
	f := LeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})

	term := RelativeMeanImprovementCrit{NumItersToAvg: 10}
	v, coords := Minimize(f, vector.Zeroes(2), 1e-10, &term, nil)
	t.Log("Value: ", v, "Coords: ", coords)

	if math.Abs(coords[0]) > 1e-2 {
		t.Error("coords[0] != 0: %s", coords[0])
	}
	if math.Abs(coords[1]-1) > 1e-2 {
		t.Error("coords[1] != 1: %s", coords[0])
	}
}

func TestLeastSquares(t *testing.T) {
	// Demo problem from https://en.wikipedia.org/wiki/Linear_least_squares_(mathematics)
	f := LeastSquares([]vector.F64{
		vector.F64{1, 6},
		vector.F64{2, 5},
		vector.F64{3, 7},
		vector.F64{4, 10},
	})

	term := RelativeMeanImprovementCrit{NumItersToAvg: 10}
	v, coords := Minimize(f, vector.Zeroes(2), 1e-8, &term, nil)
	t.Log("Value: ", v, "Coords: ", coords)
	if math.Abs(coords[0]-3.5) > 1e-1 {
		t.Error("coords[0] != 3.5: %s", coords[0])
	}
	if math.Abs(coords[1]-1.4) > 1e-1 {
		t.Error("coords[1] != 1.4: %s", coords[0])
	}
}

func BenchmarkLeastSquare(b *testing.B) {
	f := LeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})

	term := NumIterationsCrit{NumIterations: 1}

	for i := 0; i < b.N; i++ {
		Minimize(f, vector.Zeroes(2), 1e-8, &term, nil)
	}
}

func init() {
	rand.Seed(1)
}
