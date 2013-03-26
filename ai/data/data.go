// Data representation.
package data

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"math/rand"
	"reflect"
)

type Instance interface {
	View(attrs []Attr) vector.F64
	Get(attr Attr) float64
}

type Instances interface {
	Len() int
	Attrs() Attributes

	Get(idx int) Instance
	View(attrs []Attr) []vector.F64

	Perm(perm []int) Instances
	Shuffled() Instances
	Split(idx int) (Instances, Instances)
}

type arrayInstances struct {
	attrs  Attributes
	values []vector.F64 // each Values[i] is an Instance
}

func (a *arrayInstances) Len() int {
	return len(a.values)
}

func (a *arrayInstances) Attrs() Attributes {
	return a.attrs
}

func (a *arrayInstances) Shuffled() Instances {
	return a.Perm(rand.Perm(a.Len()))
}

func (a *arrayInstances) Split(idx int) (Instances, Instances) {
	values1 := make([]vector.F64, idx)
	values2 := make([]vector.F64, a.Len()-idx)

	copy(values1, a.values[:idx])
	copy(values2, a.values[idx:])
	return &arrayInstances{attrs: a.attrs, values: values1},
		&arrayInstances{attrs: a.attrs, values: values2}
}

func (a *arrayInstances) Perm(perm []int) Instances {
	values := make([]vector.F64, a.Len())

	for i, j := range perm {
		values[i] = a.values[j]
	}
	return &arrayInstances{attrs: a.attrs, values: values}
}

func (a *arrayInstances) View(attrs []Attr) []vector.F64 {
	// fast path first
	if (a.attrs.Eq(attrs)) {
		return a.values
	}

	panic("attr remapping is not implemented")
}

type arrayInstance struct {
	ai *arrayInstances
	idx int
}

func (a *arrayInstances) Get(idx int) Instance {
	return &arrayInstance{ai: a, idx: idx}
}

func (a *arrayInstance) View(attrs []Attr) vector.F64 {
	// fast path first
	if (a.ai.attrs.Eq(attrs)) {
		return a.ai.values[a.idx]
	}

	panic("attr remapping is not implemented")
}

func (a *arrayInstance) Get(attr Attr) float64 {
	return a.ai.values[a.idx][a.ai.attrs.IndexOf(attr)]
}


func Of(arr interface{}) Instances {
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

	instances := arrayInstances{
		attrs:  attrs,
		values: make([]vector.F64, length),
	}

	for i := 0; i < length; i++ {
		element := v.Index(i)
		instances.values[i] = make([]float64, numFields)

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
				instances.values[i][j] = float64(idx)

			case reflect.Bool:
				value := element.Field(j).Bool()
				if value {
					instances.values[i][j] = 1
				} else {
					instances.values[i][j] = 0
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

	return &instances
}

type zippedInstances struct {
	attrs Attributes
	ii    []Instances
}

func (zi *zippedInstances) Len() int {
	return zi.ii[0].Len()
}

func (zi *zippedInstances) Attrs() Attributes {
	return zi.attrs
}

func (zi *zippedInstances) Shuffled() Instances {
	perm := rand.Perm(zi.Len())
	return zi.Perm(perm)
}

func (zi *zippedInstances) Split(idx int) (Instances, Instances) {
	ii1 := make([]Instances, len(zi.ii))
	ii2 := make([]Instances, len(zi.ii))

	for i, instances := range zi.ii {
		ii1[i], ii2[i] = instances.Split(idx)
	}
	return &zippedInstances{attrs: zi.attrs, ii: ii1},
		&zippedInstances{attrs: zi.attrs, ii: ii2}
}

func (zi *zippedInstances) Perm(perm []int) Instances {
	ii := make([]Instances, len(zi.ii))
	for i, instances := range zi.ii {
		ii[i] = instances.Perm(perm)
	}
	return &zippedInstances{
		attrs: zi.attrs,
		ii:    ii,
	}
}

func (zi *zippedInstances) View(attrs []Attr) []vector.F64 {
	// fast path first
	for _, instances := range zi.ii {
		if (instances.Attrs().ContainsAll(attrs)) {
			return instances.View(attrs)
		}
	}

	panic("Cross-subinstances view not implemented yet.")
}

func concatAttributes(ii []Instances) Attributes {
	numAttrs := 0
	for _, instances := range ii {
		numAttrs += len(instances.Attrs())
	}

	result := Attributes(make([]Attr, numAttrs))

	index := 0
	for _, instances := range ii {
		copy(result[index:], instances.Attrs())
		index += len(instances.Attrs())
	}
	return result
}

func Zip(iinst ...Instances) Instances {
	length := iinst[0].Len()
	for _, instances := range iinst {
		if instances.Len() != length {
			panic(fmt.Errorf("Length mismatch"))
		}
	}

	return &zippedInstances{attrs: concatAttributes(iinst), ii: iinst}
}

func (zi *zippedInstances) Get(idx int) Instance {
	return &zippedInstance{zi: zi, idx: idx}
}

func FromRows(rows []vector.F64, attrs []Attr) Instances {
	return &arrayInstances{attrs: attrs, values: rows}
}

type zippedInstance struct {
	zi *zippedInstances
	idx int
}

func (zi *zippedInstance) View(attrs []Attr) vector.F64 {
	// fast path
	for _, instances := range zi.zi.ii {
		if (instances.Attrs().ContainsAll(attrs)) {
			return instances.Get(zi.idx).View(attrs)
		}
	}

	panic("cross-instance not implemented yet.")
}

func (zi *zippedInstance) Get(attr Attr) float64 {
	// fast path
	for _, instances := range zi.zi.ii {
		if (instances.Attrs().Contains(attr)) {
			return instances.Get(zi.idx).Get(attr)
		}
	}

	panic("cross-instance not implemented yet.")
}
