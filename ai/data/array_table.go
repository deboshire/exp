package data

import (
	"github.com/deboshire/exp/math/vector"
	"math/rand"
)

type arrayTable struct {
	attrs  Attributes
	values []vector.F64 // each values[i] is a row
}

type arrayRow struct {
	table *arrayTable
	idx   int
}

func FromRows(rows []vector.F64, attrs []Attr) Table {
	return &arrayTable{attrs: attrs, values: rows}
}

func (a *arrayTable) Len() int {
	return len(a.values)
}

func (a *arrayTable) Attrs() Attributes {
	return a.attrs
}

func (a *arrayTable) Shuffled() Table {
	return a.Perm(rand.Perm(a.Len()))
}

func (a *arrayTable) Split(idx int) (Table, Table) {
	values1 := make([]vector.F64, idx)
	values2 := make([]vector.F64, a.Len()-idx)

	copy(values1, a.values[:idx])
	copy(values2, a.values[idx:])
	return &arrayTable{attrs: a.attrs, values: values1},
		&arrayTable{attrs: a.attrs, values: values2}
}

func (a *arrayTable) Perm(perm []int) Table {
	values := make([]vector.F64, a.Len())

	for i, j := range perm {
		values[i] = a.values[j]
	}
	return &arrayTable{attrs: a.attrs, values: values}
}

func (a *arrayTable) View(attrs []Attr) []vector.F64 {
	// fast path first
	if a.attrs.Eq(attrs) {
		return a.values
	}

	panic("attr remapping is not implemented")
}

func (t *arrayTable) Get(idx int) Row {
	return &arrayRow{table: t, idx: idx}
}

func (t *arrayTable) idxMap(attrs []Attributes) [][]int {
	result := make([][]int, len(attrs))

	for i := range attrs {
		result[i] = make([]int, len(t.Attrs()))

		for j := range t.attrs {
			result[i][j] = attrs[i].IndexOf(t.attrs[j])
		}
	}

	return result
}

type arrayTableIterator struct {
	t      *arrayTable
	idx    int
	result []vector.F64
	idxMap [][]int
}

func (it *arrayTableIterator) next() (row []vector.F64, ok bool) {
	if it.idx >= len(it.t.values) {
		return nil, false
	}

	row = it.result

	for i := range it.result {
		for j := range it.result[i] {
			col := it.idxMap[i][j]
			if col >= 0 {
				v := it.t.values[it.idx][col]
				row[i][j] = v
			}
		}
	}

	it.idx++
	return row, true
}

func (t *arrayTable) Iterator(attrs []Attributes) Iterator {
	result := make([]vector.F64, len(attrs))
	idxMap := make([][]int, len(attrs))

	for i := range result {
		result[i] = vector.NaN(len(attrs[i]))

		idxMap[i] = make([]int, len(attrs[i]))
		for j := range idxMap[i] {
			idxMap[i][j] = t.Attrs().IndexOf(attrs[i][j])
		}
	}

	return Iterator((&arrayTableIterator{t: t, result: result, idxMap: idxMap}).next)
}

func (r *arrayRow) View(attrs []Attr) vector.F64 {
	// fast path first
	if r.table.attrs.Eq(attrs) {
		return r.table.values[r.idx]
	}

	panic("attr remapping is not implemented")
}

func (r *arrayRow) Get(attr Attr) float64 {
	return r.table.values[r.idx][r.table.attrs.IndexOf(attr)]
}
