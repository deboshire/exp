package vector

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"unsafe"
)

// #cgo CFLAGS:-O3 -ffast-math
// #include "vector.h"
import "C"

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

func (v F64) CopyTo(target F64) {
	copy(target, v)
}

func (v F64) CopyFrom(src F64) {
	copy(v, src)
}

func Zeroes(size int) F64 {
	return F64(make([]float64, size))
}

func NaN(size int) F64 {
	result := F64(make([]float64, size))
	result.Fill(math.NaN())
	return result
}

func (v F64) Fill(a float64) F64 {
	for i := range v {
		v[i] = a
	}
	return v
}

func (v F64) Sub(v1 F64) {
	assertSameLen(v, v1)
	for i := range v {
		v[i] -= v1[i]
	}
}

func (v F64) SubTo(v1 F64, dest F64) {
	assertSameLen(v, v1)
	assertSameLen(v, dest)
	for i := range v {
		dest[i] = v[i] - v1[i]
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

func addr(v F64) unsafe.Pointer {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	return unsafe.Pointer(header.Data)
}

func (v F64) Dist2(v1 F64) float64 {
	assertSameLen(v, v1)
	return float64(C.dist2(addr(v), addr(v1), C.int(len(v))))

	// d := 0.0
	// for i := range v {
	// 	a := v[i] - v1[i]
	// 	d += a * a
	// }
	// return d
}

func (v F64) Dist(v1 F64) float64 {
	return math.Sqrt(v.Dist2(v1))
}

func (v F64) Length() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v F64) Len() int {
	return len(v)
}

func (v F64) DotProduct(v1 F64) float64 {
	assertSameLen(v, v1)
	return float64(C.dot(addr(v), addr(v1), C.int(len(v))))

	// result := 0.0
	// for i := range v {
	// 	result += v[i] * v1[i]
	// }
	// return result
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

func (v F64) Eq(v1 F64, eps float64) bool {
	if len(v) != len(v1) {
		return false
	}

	for i := range v {
		if math.Abs(v[i]-v1[i]) > eps {
			return false
		}
	}
	return true
}

func (v B) Shuffle() {
	for i := len(v) - 1; i > 0; i-- {
		j := rand.Intn(i)
		v[i], v[j] = v[j], v[i]
	}
}

func Parse(strs []string) (res F64, err error) {
	res = Zeroes(len(strs))

	for i, str := range strs {
		parsedFloat, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return res, err
		}
		res[i] = parsedFloat
	}
	return
}
