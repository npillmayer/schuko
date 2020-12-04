package gotestingadapter_test

import (
	"testing"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/gconf"
	"github.com/npillmayer/schuko/testadapter"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/gotestingadapter"
)

func Test1(t *testing.T) {
	l := gotestingadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("Hello 1")
	l.P("a", "b").Infof("World")
	l.Debugf("Hello 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Debugf("Hello 3")
	// will produce no output
}

func Test2(t *testing.T) {
	schuko.AddTraceAdapter("test", gotestingadapter.GetAdapter())
	c := testadapter.New()
	c.Set("tracing", "test")
	gconf.Initialize(c)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	tracing.EngineTracer.P("key", "value").Errorf("This is a test")
	// output only seen with -v flag turned on
}
