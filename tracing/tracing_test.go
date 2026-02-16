package tracing

import (
	"bytes"
	"io"
	"testing"
)

func TestSelectorNoOp(t *testing.T) {
	tracer := Select("root")
	buf := &bytes.Buffer{} // log destination
	tracer.SetOutput(buf)
	Errorf("this is a test error")
	traceout := buf.String()
	t.Logf("trace: %q", traceout)
	if traceout != "" {
		t.Errorf("expected default Selector be a void selector; isn't")
	}
}

func TestSelectorInstallation(t *testing.T) {
	buf := &bytes.Buffer{} // log destination
	sel := &testTracer{}   // a testTracer is able to log simple error messages
	sel.SetOutput(buf)
	SetTraceSelector(sel) // install it as global selector
	msg := "this is a test error"
	Errorf(msg)              // this should log to global tracer "root"
	traceout := buf.String() // collect the output
	t.Logf("trace: %q", traceout)
	if traceout != msg {
		t.Errorf("expected root-tracer to log error %q; isn't", msg)
	}
}

// ---------------------------------------------------------------------------

type testTracer struct {
	*noOpTrace
	out io.Writer
}

func (tt *testTracer) Errorf(msg string, args ...any) {
	tt.out.Write([]byte(msg)) // for test: ignore args
}

func (tt *testTracer) SetOutput(w io.Writer) {
	tt.out = w
}

func (tt *testTracer) Select(string) Trace { // testTracer is its own selector
	return tt
}
