package vector

import (
	"fmt"
)

type F64 []float64
type B []bool

func assertSameLen(v1 F64, v2 F64) {
	if v1.Len() != v2.Len() {
		panic(fmt.Sprintf("Length mismatch: %d != %d", v1.Len(), v2.Len()))
	}
}

func (v F64) Copy() F64 {
	result := Zeroes(len(v))
	copy(result, v)
	return result
}

func Zeroes(size int) F64 {
	return F64(make([]float64, size))
}

func (v F64) Sub(v1 F64) {
	assertSameLen(v, v1)
	for i := range v {
		v[i] -= v1[i]
	}
}

func (v F64) Add(v1 F64) {
	assertSameLen(v, v1)
	for i := range v {
		v[i] += v1[i]
	}
}

func (v F64) Mul(s float64) {
	for i := range v {
		v[i] *= s
	}
}

func (v F64) Dist2(v1 F64) float64 {
	assertSameLen(v, v1)
	d := 0.0
	for i := range v {
		a := v[i] - v1[i]
		d += a * a
	}
	return d
}

func (v F64) Len() int {
	return len(v)
}

func (v F64) DotProduct(v1 F64) float64 {
	assertSameLen(v, v1)

	result := 0.0
	for i := range v {
		result += v[i] * v1[i]
	}
	return result
}

func (v F64) F64ToB() B {
	result := make([]bool, v.Len())
	for i, f := range v {
		if f >= 0.5 {
			result[i] = true
		}
	}
	return B(result)
}
