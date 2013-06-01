// Data representation.
package data

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"reflect"
)

// todo: change iterator type to better fit into for() statement.
// like:
//     for it.next(&data) { ... }
type Iterator func() (row []vector.F64, ok bool)

type Table interface {
	Len() int
	Attrs() Attributes

	Perm(perm []int) Table
	Shuffled() Table
	Split(idx int) (Table, Table)

	Iterator(attrs []Attributes) Iterator
	// Iterator that never stops.
	CyclicIterator(attrs []Attributes) Iterator

	TransformAttr(attr Attr, transform AttrTransform)

	Do(f func(row []vector.F64), attrs []Attributes)
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
