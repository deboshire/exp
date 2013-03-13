/*
	Tracer for mathematical computations.
	Like logging, but oriented towards capturing values.
*/
package tracer

import (
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"io"
	"os"
)

type Tracer interface {
	TraceFloat64(label string, value float64)
	TraceF64(label string, value vector.F64)
	TraceInt(label string, value int)
}

type tracerImpl struct {
	p string
	w io.Writer
}

func (t tracerImpl) TraceFloat64(label string, value float64) {
	fmt.Fprintf(t.w, "%s%s : %v\n", t.p, label, value)
}

func (t tracerImpl) TraceInt(label string, value int) {
	fmt.Fprintf(t.w, "%s%s : %v\n", t.p, label, value)
}

func (t tracerImpl) TraceF64(label string, value vector.F64) {
	fmt.Fprintf(t.w, "%s%s : %v\n", t.p, label, value)
}

type nullTracer struct{}

func (t nullTracer) TraceFloat64(label string, value float64) {
}

func (t nullTracer) TraceInt(label string, value int) {
}

func (t nullTracer) TraceF64(label string, value vector.F64) {
}

func NewNullTracer() Tracer {
	return nullTracer{}
}

func NewTracer(prefix string, writer io.Writer) Tracer {
	return tracerImpl{p: prefix, w: writer}
}

func NewStderrTracer(prefix string) Tracer {
	return tracerImpl{p: prefix, w: os.Stderr}
}

func DefaultTracer() Tracer {
	return NewNullTracer()
}
