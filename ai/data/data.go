// Data representation.
package data

import (
	"fmt"
	"reflect"
)

type AttributeType int

const (
	TYPE_NOMINAL = AttributeType(0) // finite set of unordered values
	TYPE_NUMERIC = AttributeType(1) // numberic value
)

type Schema struct {
	NumAttributes int
	Types         []AttributeType // attribute types
	Names         []string        // attribute values
	NumValues     []int           // 0 when not applicable
}

// type Instance struct {
// 	Schema *Schema
// 	Values []float64
// }

type Instances struct {
	Schema *Schema
	Len    int         // number of instances
	Values [][]float64 // each Values[i] is an Instance
}

func Of(arr interface{}) *Instances {
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

	schema := Schema{
		NumAttributes: numFields,
		Types:         make([]AttributeType, numFields),
		Names:         make([]string, numFields),
		NumValues:     make([]int, numFields),
	}

	fields := make([]reflect.StructField, elementType.NumField())
	for i := 0; i < numFields; i++ {
		fields[i] = elementType.Field(i)
		schema.Names[i] = fields[i].Name

		switch fields[i].Type.Kind() {
		default:
			panic(fmt.Errorf("Struct field type %v not supported", fields[i].Type.Kind()))
		case reflect.String:
			schema.Types[i] = TYPE_NOMINAL
		case reflect.Bool:
			schema.Types[i] = TYPE_NOMINAL
			schema.NumValues[i] = 2
		}
	}

	uniqueValues := make([]map[string]int, numFields)
	for i := range uniqueValues {
		uniqueValues[i] = make(map[string]int)
	}

	instances := Instances{
		Schema: &schema,
		Len:    length,
		Values: make([][]float64, length),
	}

	for i := 0; i < length; i++ {
		element := v.Index(i)
		instances.Values[i] = make([]float64, numFields)

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
				instances.Values[i][j] = float64(idx)

			case reflect.Bool:
				value := element.Field(j).Bool()
				if value {
					instances.Values[i][j] = 1
				} else {
					instances.Values[i][j] = 0
				}
			}
		}
	}

	for j, field := range fields {
		if field.Type.Kind() == reflect.String {
			schema.NumValues[j] = len(uniqueValues[j])
		}
	}

	return &instances
}
