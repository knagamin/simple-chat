package trace

import (
	"fmt"
	"io"
	"log"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func (t *tracer) Trace(a ...interface{}) {

	// [TODO] should not repeat it.
	if _, err := t.out.Write([]byte(fmt.Sprint(a...))); err != nil {
		log.Fatal("can not write by Tracer")
	}
	if _, err := t.out.Write([]byte("\n")); err != nil {
		log.Fatal("can not write by Tracer")
	}

}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}
func Off() Tracer {
	return &nilTracer{}
}
