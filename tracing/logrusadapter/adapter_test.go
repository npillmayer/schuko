package logrusadapter_test

import (
	"testing"

	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/logrusadapter"
)

func Test1(t *testing.T) {
	l := logrusadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("Hello 1")
	l.P("a", "b").Infof("World")
	l.Debugf("Hello 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Debugf("Hello 3")
}
