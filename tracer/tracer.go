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

	Sub(label string) Tracer
	Algorithm(label string) Tracer
	Iter(i int64) Tracer
	LastIter(i int64) Tracer
}

var (
	defaultTracer Tracer
)

type tracerImpl struct {
	p string
	w io.Writer
}

func (t *tracerImpl) TraceFloat64(label string, value float64) {
	fmt.Fprintf(t.w, "%s.%s : %v\n", t.p, label, value)
}

func (t *tracerImpl) TraceInt(label string, value int) {
	fmt.Fprintf(t.w, "%s.%s : %v\n", t.p, label, value)
}

func (t *tracerImpl) TraceF64(label string, value vector.F64) {
	fmt.Fprintf(t.w, "%s.%s : %v\n", t.p, label, value)
}

func (t *tracerImpl) Sub(label string) Tracer {
	if len(t.p) > 0 {
		label = "." + label
	}
	return &tracerImpl{p: t.p + label, w: t.w}
}

func (t *tracerImpl) Algorithm(name string) Tracer {
	return t.Sub(name)
}

func (t *tracerImpl) Iter(i int64) Tracer {
	return t
}

func (t *tracerImpl) LastIter(i int64) Tracer {
	return t.Iter(i)
}

type nullTracer struct{}

func (t *nullTracer) TraceFloat64(label string, value float64) {
}

func (t *nullTracer) TraceInt(label string, value int) {
}

func (t *nullTracer) TraceF64(label string, value vector.F64) {
}

func (t *nullTracer) Sub(label string) Tracer {
	return t
}

func (t *nullTracer) Algorithm(name string) Tracer {
	return t
}

func (t *nullTracer) Iter(i int64) Tracer {
	return t
}

func (t *nullTracer) LastIter(i int64) Tracer {
	return t.Iter(i)
}

func NewNullTracer() Tracer {
	return &nullTracer{}
}

func NewTracer(prefix string, writer io.Writer) Tracer {
	return &tracerImpl{p: prefix, w: writer}
}

func NewStderrTracer(prefix string) Tracer {
	return &tracerImpl{p: prefix, w: os.Stderr}
}

func NewFileTracer(fileName string) Tracer {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return &tracerImpl{p: "", w: f}
}

func DefaultTracer() Tracer {
	return defaultTracer
}

func SetDefaultTracer(t Tracer) {
	defaultTracer = t
}

func init() {
	defaultTracer = NewNullTracer()
}
