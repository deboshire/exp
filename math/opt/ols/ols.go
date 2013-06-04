// Ordinart Least Squares
package ols

import (
	"github.com/deboshire/exp/math/opt/sgd"
	"github.com/deboshire/exp/math/vector"
	"math/rand"
)

/*
	Objective function for performing least squares optimization.
*/
func SgdLeastSquares(points []vector.F64) *sgd.Minimizer {
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

	return &sgd.Minimizer{F: f, X0: vector.Zeroes(dim)}
}
