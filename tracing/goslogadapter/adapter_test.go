package goslogadapter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/goslogadapter"
)

func TestLevelFiltering(t *testing.T) {
	l := goslogadapter.New()
	buf := &bytes.Buffer{}
	l.SetOutput(buf)
	l.SetTraceLevel(tracing.LevelError)

	l.Debugf("debug msg")
	l.Infof("info msg")
	l.Errorf("error msg")

	out := buf.String()
	if strings.Contains(out, "debug msg") {
		t.Errorf("expected debug message to be filtered, got %q", out)
	}
	if strings.Contains(out, "info msg") {
		t.Errorf("expected info message to be filtered, got %q", out)
	}
	if !strings.Contains(out, "error msg") {
		t.Errorf("expected error message to be logged, got %q", out)
	}
}

func TestPFieldsAndChaining(t *testing.T) {
	l := goslogadapter.New()
	buf := &bytes.Buffer{}
	l.SetOutput(buf)
	l.SetTraceLevel(tracing.LevelDebug)

	l.P("a", "b").P("x", 7).Infof("hello %s", "world")
	out := buf.String()

	if !strings.Contains(out, "hello world") {
		t.Errorf("expected formatted message in output, got %q", out)
	}
	if !strings.Contains(out, "a=b") {
		t.Errorf("expected field a=b in output, got %q", out)
	}
	if !strings.Contains(out, "x=7") {
		t.Errorf("expected field x=7 in output, got %q", out)
	}
}

func TestTraceLevelRoundtrip(t *testing.T) {
	l := goslogadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	if lvl := l.GetTraceLevel(); lvl != tracing.LevelDebug {
		t.Errorf("expected level debug, got %v", lvl)
	}
	l.SetTraceLevel(tracing.LevelInfo)
	if lvl := l.GetTraceLevel(); lvl != tracing.LevelInfo {
		t.Errorf("expected level info, got %v", lvl)
	}
	l.SetTraceLevel(tracing.LevelError)
	if lvl := l.GetTraceLevel(); lvl != tracing.LevelError {
		t.Errorf("expected level error, got %v", lvl)
	}
}

func TestSetOutputRebindsWriter(t *testing.T) {
	l := goslogadapter.New()
	l.SetTraceLevel(tracing.LevelInfo)

	buf1 := &bytes.Buffer{}
	l.SetOutput(buf1)
	l.Infof("first")

	buf2 := &bytes.Buffer{}
	l.SetOutput(buf2)
	l.Infof("second")

	out1 := buf1.String()
	out2 := buf2.String()

	if !strings.Contains(out1, "first") {
		t.Errorf("expected first message in first writer, got %q", out1)
	}
	if strings.Contains(out1, "second") {
		t.Errorf("did not expect second message in first writer, got %q", out1)
	}
	if !strings.Contains(out2, "second") {
		t.Errorf("expected second message in second writer, got %q", out2)
	}
}

func TestGetAdapter(t *testing.T) {
	a := goslogadapter.GetAdapter()
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a() == nil {
		t.Fatal("expected adapter to return tracer")
	}
}

