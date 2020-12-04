package gologadapter_test

import (
	"testing"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/testadapter"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/gologadapter"
)

func Test1(t *testing.T) {
	l := gologadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("Hello 1")
	l.P("a", "b").Infof("World")
	l.Debugf("Hello 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Debugf("Hello 3")
}

func Test2(t *testing.T) {
	config.Initialize(testadapter.New())
	tracing.EngineTracer.P("key", "value").Errorf("This is a test")
}

func Test3(t *testing.T) {
	v := []int{1, 2, 3}
	tracing.With(tracing.EngineTracer).Dump("v", v)
}
