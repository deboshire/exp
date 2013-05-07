package data

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"math"
)

type AttrKind int16

const (
	KIND_NOMINAL = AttrKind(0) // finite set of unordered values
	KIND_NUMERIC = AttrKind(1) // numberic value
)

type AttrType struct {
	Kind      AttrKind
	NumValues int16
}

var (
	TYPE_BOOL    = AttrType{Kind: KIND_NOMINAL, NumValues: 2}
	TYPE_NUMERIC = AttrType{Kind: KIND_NUMERIC}
)

type Attr struct {
	Name string
	Type AttrType
}

type Attributes []Attr

type AttrTransform interface {
	transform(attr Attr, attrs Attributes, values []vector.F64)
}

type toNominal struct {
}

var (
	TO_NOMINAL AttrTransform = &toNominal{}
)

func (attrs Attributes) ByName(name string) Attr {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr
		}
	}

	panic(fmt.Errorf("Attribute %s not found in %v", name, attrs))
}

func (attrs Attributes) Without(attr Attr) []Attr {
	result := make([]Attr, len(attrs)-1)
	found := false

	idx := 0
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Name != attr.Name {
			result[idx] = attrs[i]
			idx++
		} else {
			found = true
		}
	}

	if !found {
		panic(fmt.Errorf("Attribute %v not found", attr))
	}
	return result
}

func (attrs Attributes) Contains(attr Attr) bool {
	for _, attr1 := range attrs {
		if attr1 == attr {
			return true
		}
	}
	return false
}

func (attrs Attributes) ContainsAll(arr []Attr) bool {
	for _, attr1 := range arr {
		found := false
		for _, attr2 := range attrs {
			if attr2 == attr1 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (attrs Attributes) Eq(arr []Attr) bool {
	if len(attrs) != len(arr) {
		return false
	}

	for i, attr1 := range attrs {
		if arr[i] != attr1 {
			return false
		}
	}

	return true
}

func (attrs Attributes) IndexOf(attr Attr) int {
	for i, attr1 := range attrs {
		if attr == attr1 {
			return i
		}
	}

	return -1
}

func (attr Attr) Repeat(n int) []Attr {
	result := make([]Attr, n)
	for i := range result {
		result[i] = attr
		result[i].Name = fmt.Sprintf("%s.%d", attr.Name, i)
	}
	return result
}

func (t toNominal) transform(attr Attr, attrs Attributes, values []vector.F64) {
	maxValue := 0
	idx := attrs.IndexOf(attr)

	for _, row := range values {
		i := int(row[idx])
		row[idx] = float64(i)
		if i > maxValue {
			maxValue = i
		}
	}

	if maxValue > math.MaxInt16 {
		panic("too many values")
	}

	attrs[idx] = Attr{Name: attr.Name, Type: AttrType{Kind: KIND_NOMINAL, NumValues: int16(maxValue + 1)}}
}
