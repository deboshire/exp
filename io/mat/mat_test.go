package mat

import (
	"bufio"
//	"fmt"
	"os"
	"testing"
	"reflect"
)

func TestFail(t *testing.T) {
	var file *os.File
	var err error

	if file, err = os.Open(os.ExpandEnv("Train1X.mat")); err != nil {
		t.Fatal(err)
	}

	d, err := Read(bufio.NewReader(file))
	
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
	
	if arr.name != "Train1X" {
		t.Fatalf("Bad name: %s", arr.name)
	}

	if len(arr.dimensions) != 2 || arr.dimensions[0] != 200 || arr.dimensions[1] != 129 {
		t.Fatalf("Bad dims: %s", arr.dimensions)
	}
	
	if len(arr.data) != 200 * 129 {
		t.Fatalf("Bad data size: %s", len(arr.data))
	}
}
