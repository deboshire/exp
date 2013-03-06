/*
	Tracer for mathematical computations.
	Like logging, but oriented towards capturing values.
*/
package tracer

import (
	"io"
	"fmt"
	"github.com/deboshire/exp/math/vector"
)

type Tracer interface {
	TraceFloat64(label string, value float64)
	TraceV64(label string, value vector.V64)
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

func (t tracerImpl) TraceV64(label string, value vector.V64) {
	fmt.Fprintf(t.w, "%s%s : %v\n", t.p, label, value) 
}

func NewTracer(prefix string, writer io.Writer) Tracer {
	return tracerImpl{p: prefix, w: writer}
}

type nullTracer struct { }


func (t nullTracer) TraceFloat64(label string, value float64) {
}

func (t nullTracer) TraceInt(label string, value int) {
}

func (t nullTracer) TraceV64(label string, value vector.V64) {
}

func NewNullTracer() Tracer {
	return nullTracer{}
}

func DefaultTracer() Tracer {
	return NewNullTracer()
}