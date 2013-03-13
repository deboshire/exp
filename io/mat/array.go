// Array manipulation routines
package mat

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
)

type Array struct {
	Name string
	Dim  []int32
	Data []float64
}

func (a *Array) RowsToVectors() (vectors []vector.F64) {
	if len(a.Dim) != 2 {
		panic(fmt.Sprintf("Array is not 2-dimensional: %s", a.Dim))
	}

	rows := int(a.Dim[0])
	rowLen := int(a.Dim[1])
	data := a.Data

	for i := 0; i < rows; i++ {
		row := data[(i * rowLen):((i + 1) * rowLen)]
		vectors = append(vectors, row)
	}

	return vectors
}

func (a *Array) ToVector() vector.F64 {
	if len(a.Dim) != 2 {
		panic(fmt.Sprintf("Array is not 2-dimensional: %s", a.Dim))
	}

	if a.Dim[0] != 1 && a.Dim[1] != 1 {
		panic(fmt.Sprintf("One of dimensions should be 1: %s", a.Dim))
	}

	return vector.F64(a.Data)
}
