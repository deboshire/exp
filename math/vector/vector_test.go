package vector

import (
	"testing"
)

func TestDist2(t *testing.T) {
	v1 := F64{1, 2}
	v2 := F64{2, 3}

	if v1.Dist2(v2) != 2 {
		t.Error("Bad distance:", v1.Dist2(v2))
	}
}


func BenchmarkDist2(b *testing.B) {
	v1 := Zeroes(10000)
	v2 := Zeroes(10000)

	for i := 0; i < b.N; i++ {
		v1.Dist2(v2)
	}
}

func BenchmarkDotProduct(b *testing.B) {
	v1 := Zeroes(10000).Fill(1)
	v2 := Zeroes(10000).Fill(2)

	for i := 0; i < b.N; i++ {
		v1.DotProduct(v2)
	}
}
