package mat

import (
	"bufio"
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

	d, err := read0(bufio.NewReader(file))

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
}
