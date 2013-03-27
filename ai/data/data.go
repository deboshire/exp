// Data representation.
package data

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"math/rand"
	"reflect"
)

type Row interface {
	View(attrs []Attr) vector.F64
	Get(attr Attr) float64
}

type Table interface {
	Len() int
	Attrs() Attributes

	Get(idx int) Row
	View(attrs []Attr) []vector.F64

	Perm(perm []int) Table
	Shuffled() Table
	Split(idx int) (Table, Table)
}

type arrayTable struct {
	attrs  Attributes
	values []vector.F64 // each values[i] is a row
}

type arrayRow struct {
	table *arrayTable
	idx   int
}

type zippedTable struct {
	attrs  Attributes
	tables []Table
}

type zippedRow struct {
	table *zippedTable
	idx   int
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

func Of(arr interface{}) Table {
	t := reflect.TypeOf(arr)

	if t.Kind() != reflect.Slice {
		panic(fmt.Errorf("Type %v not supported", t))
	}

	v := reflect.ValueOf(arr)
	length := v.Len()
	element := v.Index(0)
	elementType := element.Type()

	if elementType.Kind() != reflect.Struct {
		panic(fmt.Errorf("Individual items should be structs; %v found", elementType))
	}

	numFields := elementType.NumField()

	attrs := make([]Attr, numFields)

	fields := make([]reflect.StructField, elementType.NumField())
	for i := 0; i < numFields; i++ {
		fields[i] = elementType.Field(i)
	}

	uniqueValues := make([]map[string]int, numFields)
	for i := range uniqueValues {
		uniqueValues[i] = make(map[string]int)
	}

	table := arrayTable{
		attrs:  attrs,
		values: make([]vector.F64, length),
	}

	for i := 0; i < length; i++ {
		element := v.Index(i)
		table.values[i] = make([]float64, numFields)

		for j := range fields {
			switch fields[j].Type.Kind() {
			case reflect.String:
				fieldValues := uniqueValues[j]
				value := element.Field(j).String()
				idx, ok := fieldValues[value]
				if !ok {
					idx = len(fieldValues)
					fieldValues[value] = idx
				}
				table.values[i][j] = float64(idx)

			case reflect.Bool:
				value := element.Field(j).Bool()
				if value {
					table.values[i][j] = 1
				} else {
					table.values[i][j] = 0
				}
			}
		}
	}

	for i := 0; i < numFields; i++ {
		fields[i] = elementType.Field(i)
		attrs[i].Name = fields[i].Name

		switch fields[i].Type.Kind() {
		default:
			panic(fmt.Errorf("Struct field type %v not supported", fields[i].Type.Kind()))
		case reflect.String:
			attrs[i].Type = AttrType{Kind: KIND_NOMINAL, NumValues: int16(len(uniqueValues[i]))}
		case reflect.Bool:
			attrs[i].Type = AttrType{Kind: KIND_NOMINAL, NumValues: 2}
		}
	}

	return &table
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
		attrs: t.attrs,
		tables:     tables,
	}
}

func (t *zippedTable) View(attrs []Attr) []vector.F64 {
	// fast path first
	for _, table := range t.tables {
		if table.Attrs().ContainsAll(attrs) {
			return table.View(attrs)
		}
	}

	panic("Cross-subtable view not implemented yet.")
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

func Zip(tables ...Table) Table {
	length := tables[0].Len()
	for _, table := range tables {
		if table.Len() != length {
			panic(fmt.Errorf("Length mismatch"))
		}
	}

	return &zippedTable{attrs: concatAttributes(tables), tables: tables}
}

func (t *zippedTable) Get(idx int) Row {
	return &zippedRow{table: t, idx: idx}
}

func FromRows(rows []vector.F64, attrs []Attr) Table {
	return &arrayTable{attrs: attrs, values: rows}
}

func (r *zippedRow) View(attrs []Attr) vector.F64 {
	// fast path
	for _, table := range r.table.tables {
		if table.Attrs().ContainsAll(attrs) {
			return table.Get(r.idx).View(attrs)
		}
	}

	panic("cross-row not implemented yet.")
}

func (r *zippedRow) Get(attr Attr) float64 {
	// fast path
	for _, table := range r.table.tables {
		if table.Attrs().Contains(attr) {
			return table.Get(r.idx).Get(attr)
		}
	}

	panic("cross-row not implemented yet.")
}
