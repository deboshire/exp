package data

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"math/rand"
)

type zippedTable struct {
	attrs  Attributes
	tables []Table
}

type zippedRow struct {
	table *zippedTable
	idx   int
}

func Zip(tables ...Table) Table {
	length := tables[0].Len()
	for _, table := range tables {
		if table.Len() != length {
			panic(fmt.Errorf("Length mismatch"))
		}
	}

	return &zippedTable{attrs: concatAttributes(tables), tables: tables}
}

func (t *zippedTable) Len() int {
	return t.tables[0].Len()
}

func (t *zippedTable) Attrs() Attributes {
	return t.attrs
}

func (zi *zippedTable) Shuffled() Table {
	perm := rand.Perm(zi.Len())
	return zi.Perm(perm)
}

func (t *zippedTable) Split(idx int) (Table, Table) {
	tables1 := make([]Table, len(t.tables))
	tables2 := make([]Table, len(t.tables))

	for i, table := range t.tables {
		tables1[i], tables2[i] = table.Split(idx)
	}
	return &zippedTable{attrs: t.attrs, tables: tables1},
		&zippedTable{attrs: t.attrs, tables: tables2}
}

func (t *zippedTable) Perm(perm []int) Table {
	tables := make([]Table, len(t.tables))
	for i, table := range t.tables {
		tables[i] = table.Perm(perm)
	}
	return &zippedTable{
		attrs:  t.attrs,
		tables: tables,
	}
}

func concatAttributes(ii []Table) Attributes {
	numAttrs := 0
	for _, table := range ii {
		numAttrs += len(table.Attrs())
	}

	result := Attributes(make([]Attr, numAttrs))

	index := 0
	for _, table := range ii {
		copy(result[index:], table.Attrs())
		index += len(table.Attrs())
	}
	return result
}

func (t *zippedTable) iterator(attrs []Attributes, iterators []Iterator) Iterator {
	data := make([]vector.F64, len(attrs))

	for i, a := range attrs {
		data[i] = vector.Zeroes(len(a))
	}

	// idx[attr_index][column] = iterator index
	idx := make([][]int, len(attrs))
	for i, a := range attrs {
		idx[i] = make([]int, len(a))
		for j := range idx[i] {
			idx[i][j] = -1
		}

		for j, attr := range a {
			for k, t := range t.tables {
				if t.Attrs().Contains(attr) {
					if idx[i][j] != -1 {
						panic("double source")
					} else {
						idx[i][j] = k
					}
				}
			}
		}
	}

	rows := make([][]vector.F64, len(attrs))

	return func() (row []vector.F64, ok bool) {
		notOks := 0
		for i, it := range iterators {
			rows[i], ok = it()

			if !ok {
				notOks++
			}
		}

		if notOks != 0 {
			if notOks != len(iterators) {
				panic("Iterators stopped not at the same time.")
			}

			return nil, false
		}

		for i := range attrs {
			for j := range attrs[i] {
				data[i][j] = rows[idx[i][j]][i][j]
			}
		}

		return data, true
	}
}

func (t *zippedTable) Iterator(attrs []Attributes) Iterator {
	iterators := make([]Iterator, len(t.tables))
	for i, t := range t.tables {
		iterators[i] = t.Iterator(attrs)
	}

	return t.iterator(attrs, iterators)
}

func (t *zippedTable) CyclicIterator(attrs []Attributes) Iterator {
	iterators := make([]Iterator, len(t.tables))
	for i, t := range t.tables {
		iterators[i] = t.CyclicIterator(attrs)
	}

	return t.iterator(attrs, iterators)
}

func (t *zippedTable) TransformAttr(attr Attr, transform AttrTransform) {
	panic("not implemented")
}

func (t *zippedTable) Do(f func(row []vector.F64), attrs []Attributes) {
	it := t.Iterator(attrs)

	for {
		row, ok := it()
		if !ok {
			return
		}
		f(row)
	}
}
