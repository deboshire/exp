package tracer

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"time"
)

type consoleTracer struct {
	ts      int64
	verbose bool
}

func NewConsoleTracer() Tracer {
	return &consoleTracer{}
}

func (t *consoleTracer) TraceFloat64(label string, value float64) {
	if t.verbose {
		fmt.Println(label, value)
	}
}

func (t *consoleTracer) TraceInt(label string, value int) {
	if t.verbose {
		fmt.Println(label, value)
	}
}

func (t *consoleTracer) TraceF64(label string, value vector.F64) {
	if t.verbose {
		fmt.Println(label, value)
	}
}

func (t *consoleTracer) Sub(label string) Tracer {
	return &consoleTracer{verbose: true}
}

func (t *consoleTracer) Algorithm(name string) Tracer {
	fmt.Println("# starting", name)
	return t.Sub(name)
}

func (t *consoleTracer) Iter(i int64) Tracer {
	if t.ts == 0 || time.Now().UnixNano()-t.ts > 5000000000 {
		t.ts = time.Now().UnixNano()
		t.verbose = true
		fmt.Println("---")
		fmt.Println("i", i)
	} else {
		t.verbose = false
	}
	return t
}

func (t *consoleTracer) LastIter(i int64) Tracer {
	t.verbose = true
	fmt.Println("--- last iter")
	fmt.Println("i", i)
	return t
}
