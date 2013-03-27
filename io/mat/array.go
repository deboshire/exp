// Array manipulation routines
package mat

import (
	"fmt"
	"github.com/deboshire/exp/ai/data"
	"github.com/deboshire/exp/math/vector"
)

type Array struct {
	Name string
	Dim  []int32
	Data []float64
}

func (a *Array) RowsToVectors() (vectors []vector.F64) {
	if len(a.Dim) != 2 {
		panic(fmt.Sprintf("Array is not 2-dimensional: %v", a.Dim))
	}

	rows := int(a.Dim[0])
	rowLen := int(a.Dim[1])
	data := a.Data

	for i := 0; i < rows; i++ {
		row := make([]float64, rowLen)
		for j := 0; j < rowLen; j++ {
			row[j] = data[i+j*rows]
		}
		vectors = append(vectors, row)
	}

	return
}

func (a *Array) ToVector() vector.F64 {
	if len(a.Dim) != 2 {
		panic(fmt.Sprintf("Array is not 2-dimensional: %v", a.Dim))
	}

	if a.Dim[0] != 1 && a.Dim[1] != 1 {
		panic(fmt.Sprintf("One of dimensions should be 1: %v", a.Dim))
	}

	return vector.F64(a.Data)
}

func (a *Array) RowsToTable() data.Table {
	if len(a.Dim) != 2 {
		panic(fmt.Sprintf("Array is not 2-dimensional: %v", a.Dim))
	}

	numRows := int(a.Dim[0])
	rowLen := int(a.Dim[1])
	rows := make([]vector.F64, numRows)

	for i := range rows {
		rows[i] = vector.Zeroes(rowLen)
		for j := 0; j < rowLen; j++ {
			rows[i][j] = a.Data[i+j*numRows]
		}
	}

	return data.FromRows(rows, data.Attr{Name: a.Name, Type: data.TYPE_NUMERIC}.Repeat(rowLen))
}

func (a *Array) Rename(name string) *Array {
	return &Array{Name: name, Dim: a.Dim, Data: a.Data}
}
