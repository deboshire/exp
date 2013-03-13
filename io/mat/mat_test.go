package mat

import (
	"bufio"
	v "github.com/deboshire/exp/math/vector"
	"os"
	"reflect"
	"testing"
)

func TestFail(t *testing.T) {
	var file *os.File
	var err error

	if file, err = os.Open("Train1X.mat"); err != nil {
		t.Fatal(err)
	}

	d, err := read(bufio.NewReader(file))

	if err != nil {
		t.Fatal(err)
	}

	if len(d) != 1 {
		t.Fatalf("Wrong length: %s", d)
	}

	arr, ok := d[0].(Array)
	if !ok {
		t.Fatalf("Bad element type: %s", reflect.TypeOf(d[0]))
	}

	if arr.Name != "Train1X" {
		t.Fatalf("Bad name: %s", arr.Name)
	}

	if len(arr.Dim) != 2 || arr.Dim[0] != 200 || arr.Dim[1] != 129 {
		t.Fatalf("Bad dims: %s", arr.Dim)
	}

	if len(arr.Data) != 200*129 {
		t.Fatalf("Bad data size: %s", len(arr.Data))
	}

	vectors := arr.RowsToVectors()
	if len(vectors) != 200 {
		t.Fatalf("vectors size mismatch: %d", len(vectors))
	}

	vector := vectors[0]
	if len(vector) != 129 {
		t.Fatalf("vector len mismatch: %d", len(vector))
	}

	if !vectors[0].Eq(
		v.F64{
			// cross checked with octave
			1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		1e-10) {
		t.Fatalf("vector mismatch: %#v", vectors[0])
	}
}
