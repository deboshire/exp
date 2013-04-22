package stat

import (
	"testing"
)

func Test(t *testing.T) {
	w := Window{N: 3}
	w.Add(1)
	w.Add(3)
	w.Add(2)

	if w.Max() != 3 {
		t.Errorf("Bad max: %v", w)
	}

	if w.Mean() != 2 {
		t.Errorf("Bad mean: %v %v", w, w.Mean())
	}

	w.Add(1)
	w.Add(1)

	if w.Max() != 2 {
		t.Errorf("Bad max: %v", w)
	}
	if w.Mean() != 1.3333333333333333 {
		t.Errorf("Bad mean: %v %v", w, w.Mean())
	}
}
