package trace2go_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/npillmayer/schuko/schukonf/testadapter"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/trace2go"
)

func TestRootBlank(t *testing.T) {
	root := trace2go.Root()

	buf := &bytes.Buffer{} // log destination
	root.SetOutput(buf)
	root.Infof("this is an info")
	traceout := buf.String()
	t.Logf("root trace [info]:  %q", traceout)
	if traceout != "" {
		t.Errorf("expected root tracer to be a bare bones tracer; isn't")
	}

	buf = &bytes.Buffer{} // log destination
	root.SetOutput(buf)
	msg := "this is a test error"
	root.Errorf(msg)
	traceout = buf.String()
	t.Logf("root trace [error]: %q", traceout)
	if traceout != msg {
		t.Errorf("expected root tracer to log errors; didn't")
	}
}

func TestSelectorInstallation(t *testing.T) {
	buf := &bytes.Buffer{} // log destination
	root := trace2go.Root()
	root.SetOutput(buf)
	tracing.SetTraceSelector(trace2go.Selector()) // install trace2go as global selector
	msg := "this is a test error"
	tracing.Errorf(msg)      // this should log to global tracer "root"
	traceout := buf.String() // collect the output
	t.Logf("trace: %q", traceout)
	if traceout != msg {
		t.Errorf("expected root-tracer to log errors; didn't")
	}
}

func TestSelection(t *testing.T) {
	tracing.RegisterTraceAdapter("test", getTT, true)
	tracing.SetTraceSelector(trace2go.Selector()) // install trace2go as global selector
	conf := testadapter.New()                     // lightweight configuration
	conf.Set("tracing.adapter", "test")           // test.adapter will adapt to testTracer below
	conf.Set("LEVEL.my.new.trace", "Info")        // test tracer should have level info
	trace2go.ConfigureRoot(conf, "LEVEL")         // root will spawn 'testTracer' children
	tracer := tracing.Select("my.new.trace")      // now get tracer from factory
	buf := &bytes.Buffer{}                        // log destination
	tracer.SetOutput(buf)
	msg := "this is a test info"
	tracer.Infof(msg)        // this should log to my.new.trace at Info level
	traceout := buf.String() // collect the output
	t.Logf("trace: %q", traceout)
	if traceout != msg {
		t.Errorf("expected root-tracer to log infos; didn't")
	}
}

// ---------------------------------------------------------------------------

func getTT() tracing.Trace {
	return newTestTracer()
}

func newTestTracer() *testTracer {
	tt := &testTracer{
		out: os.Stderr,
	}
	tt.Trace = tracing.NoOpTrace()
	return tt
}

type testTracer struct {
	tracing.Trace
	out io.Writer
}

func (tt *testTracer) Infof(msg string, args ...interface{}) {
	tt.out.Write([]byte(msg)) // for test: ignore args
}

func (tt *testTracer) SetOutput(w io.Writer) {
	tt.out = w
}
