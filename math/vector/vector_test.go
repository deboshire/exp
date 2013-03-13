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
