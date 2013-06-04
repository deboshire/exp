package ols

import (
	"github.com/deboshire/exp/math/opt"
	"github.com/deboshire/exp/math/vector"
	"github.com/deboshire/exp/tracer"
	"math/rand"
	"testing"
)

func TestLeastSquaresPrecise(t *testing.T) {
	minimizer := SgdLeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})
	minimizer.Tracer = tracer.DefaultTracer().Sub("TestLeastSquaresPrecise")

	term := opt.MaxRelativeChangeCrit{}
	coords := minimizer.Minimize(1e-10, &term)
	t.Log("Coords: ", coords)

	expectedCoords := vector.F64{0, 1}
	if coords.Dist(expectedCoords) > 1e-2 {
		t.Errorf("coords != expectedCoords: %v != %v", coords, expectedCoords)
	}
}

func TestLeastSquares_NumIterations(t *testing.T) {
	// Demo problem from https://en.wikipedia.org/wiki/Linear_least_squares_(mathematics)
	minimizer := SgdLeastSquares([]vector.F64{
		vector.F64{1, 6},
		vector.F64{2, 5},
		vector.F64{3, 7},
		vector.F64{4, 10},
	})
	minimizer.Tracer = tracer.DefaultTracer().Sub("TestLeastSquaresPrecise")

	term := opt.MaxRelativeChangeCrit{}
	coords := minimizer.Minimize(1e-8, &term)
	t.Log("Coords: ", coords)
	expectedCoords := vector.F64{3.5, 1.4}
	if coords.Dist(expectedCoords) > 1e-2 {
		t.Errorf("coords != expectedCoords: %v != %v", coords, expectedCoords)
	}
}

func BenchmarkLeastSquare(b *testing.B) {
	minimizer := SgdLeastSquares([]vector.F64{
		vector.F64{1, 1},
		vector.F64{2, 2},
	})

	term := opt.NumIterationsCrit{NumIterations: b.N}
	minimizer.Minimize(1e-8, &term)
}

func init() {
	rand.Seed(1)
	// tracer.SetDefaultTracer(tracer.NewConsoleTracer())
}
