package vector

type V64 []float64

func (v V64) Copy() V64 {
	result := Zeroes(len(v))
	copy(v, result)
	return result
}

func Zeroes(size int) V64 {
	return V64(make([]float64, size))
}

func (v V64) Sub(v1 V64) {
	for i := range v {
		v[i] -= v1[i]
	}
}

func (v V64) Add(v1 V64) {
	for i := range v {
		v[i] += v1[i]
	}
}

func (v V64) Mul(s float64) {
	for i := range v {
		v[i] *= s
	}
}

func (v V64) Dist2(v1 V64) float64 {
	d := 0.0
	for i := range v {
		a := v[i] - v1[i]
		d += a * a
	}
	return d
}
