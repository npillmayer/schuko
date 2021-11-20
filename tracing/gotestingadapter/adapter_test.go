package gotestingadapter_test

import (
	"testing"

	"github.com/npillmayer/schuko/schukonf/testconfig"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/gotestingadapter"
	"github.com/npillmayer/schuko/tracing/trace2go"
)

// I have not found a way to capture the output of t.Logf.
// That means, these tests are to be manually checked :-(

// output only seen with -v flag turned on
func Test1(t *testing.T) {
	l := gotestingadapter.New(t)
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("test: Test message 1")
	l.P("a", "b").Infof("test: Another test message")
	l.Debugf("test: Test message 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Errorf("test: Test message 3")
}

// output only seen with -v flag turned on
func Test2(t *testing.T) {
	tracing.RegisterTraceAdapter("test", gotestingadapter.GetAdapter(t), true)
	c := testconfig.Conf{
		"tracing.adapter": "test",
		"LEVEL.x":         "Debug",
	}
	if err := trace2go.ConfigureRoot(c, "LEVEL"); err != nil {
		t.Fatal(err)
	}
	tracing.SetTraceSelector(trace2go.Selector())
	tracer := tracing.Select("x")
	tracer.Infof("test: This is an info message")
	trace2go.Teardown()
}

func Test3(t *testing.T) {
	tracer := tracing.Select("x")
	// this goes through the default root tracer, which should not react to level info
	tracer.Infof("test: This info message should not be displayed")
	teardown := gotestingadapter.QuickConfig(t, "x")
	defer teardown()
	tracer = tracing.Select("x")
	// both messages should be displayed
	tracer.Debugf("test: This is a debug message")
	tracer.Errorf("test: This is an error message")
}
