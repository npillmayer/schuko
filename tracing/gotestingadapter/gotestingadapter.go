/*
Package gotestingadapter implements tracing with the Go testing  logger.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
gotestingadapter uses the Go testing logging mechanism, i.e. "t.logf(...)",
with t of type *testing.T.

As we are logging to global tracers there is no way of configuring them
specifically from single tests. That's not a good thing, as in general
tests may be executed concurrently/in parallel. This would confuse the
tracing.

# Attention

Clients must not use the testingtracer in concurrent mode.
Please set

	go test -p 1

# BSD License

# Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package gotestingadapter

import (
	"fmt"
	"io"
	"testing"

	"github.com/npillmayer/schuko/schukonf/testconfig"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/trace2go"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a Go testing logger.
type Tracer struct {
	t     *testing.T
	p     string
	level tracing.TraceLevel
}

var logLevelPrefix = []string{"ERROR ", "INFO  ", "DEBUG "}

//var allTracers =

// New creates a new Tracer instance valid for a testing.T.
func New(t *testing.T) tracing.Trace {
	return &Tracer{
		t:     t,
		p:     "",
		level: tracing.LevelError,
	}
}

// GetAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter(t *testing.T) tracing.Adapter {
	traceT := t
	return func() tracing.Trace {
		return &Tracer{
			t:     traceT,
			p:     "",
			level: tracing.LevelError,
		}
	}
}

// As we are logging to global tracers there is no way of configuring them
// specifically from single tests. That's not a good thing, as in general
// tests may be executed concurrently or in parallel. This would confuse the
// tracing.
// Users have to be aware that the testingtracer may not be used concurrently.
//
// Deprecated: Will be removed, not needed any more.
var globalTestingT *testing.T

// RedirectTracing will be called by clients at the start of a test. This will
// redirect all (global) tracers to use t.logf(...).
//
// It returns a teardown function which should be called at the end of a test.
// The usual pattern will look like this:
//
//	func TestSomething(t *testing.T) {
//	     teardown := gotestingadapter.RedirectTracing(t)
//	     defer teardown()
//	     ...
//	 }
//
// Deprecated: Will be removed, not needed any more.
func RedirectTracing(t *testing.T) func() {
	globalTestingT = t
	return teardownTestingT
}

// Deprecated: Will be removed, not needed any more.
func teardownTestingT() {
	globalTestingT = nil
}

// P is part of interface Trace
func (tr *Tracer) P(key string, val any) tracing.Trace {
	var p string
	switch v := val.(type) {
	case rune:
		p = fmt.Sprintf("[%s=%#U] ", key, v)
	case int, int8, int16, int64, uint16, uint32, uint64:
		p = fmt.Sprintf("[%s=%d] ", key, v)
	case string:
		p = fmt.Sprintf("[%s=%s] ", key, v)
	default:
		p = fmt.Sprintf("[%s=%v] ", key, v)
	}
	tr.p = p
	return tr
}

func (tr *Tracer) output(l tracing.TraceLevel, s string, args ...any) {
	if tr.t != nil {
		prefix := fmt.Sprintf("%s%s", logLevelPrefix[int(l)], tr.p)
		tr.t.Logf(prefix+s, args...)
		tr.p = ""
	} else if globalTestingT != nil {
		prefix := fmt.Sprintf("%s%s", logLevelPrefix[int(l)], tr.p)
		globalTestingT.Logf("depr."+prefix+s, args...)
		tr.p = ""
	}
}

// Debugf is part of interface Trace
func (tr *Tracer) Debugf(s string, args ...any) {
	if tr.level < tracing.LevelDebug {
		return
	}
	tr.output(tracing.LevelDebug, s, args...)
}

// Infof is part of interface Trace
func (tr *Tracer) Infof(s string, args ...any) {
	if tr.level < tracing.LevelInfo {
		return
	}
	tr.output(tracing.LevelInfo, s, args...)
}

// Errorf is part of interface Trace
func (tr *Tracer) Errorf(s string, args ...any) {
	if tr.level < tracing.LevelError {
		return
	}
	tr.output(tracing.LevelError, s, args...)
}

// SetTraceLevel is part of interface Trace
func (tr *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	tr.p = ""
	tr.level = l
}

// GetTraceLevel is part of interface Trace
func (tr *Tracer) GetTraceLevel() tracing.TraceLevel {
	return tr.level
}

// SetOutput is part of interface Trace. This implementation ignores it.
func (tr *Tracer) SetOutput(writer io.Writer) {}

// ----------------------------------------------------------------------

// QuickConfig sets up a configuration suitable for test cases, including tracing.
// It returns a teardown function which should be called at the end of a test.
// The usual pattern will look like this:
//
//	func TestSomething(t *testing.T) {
//	     teardown := testconfig.QuickConfig(t, "first.trace.name", "second.trace.name")
//	     defer teardown()
//	     …
//	 }
//
// Tracing output will be redirected to the testing.T log (`t.Logf(…)`).
// All tracers identified by "first.trace.name" etc. will be created and have their log
// levels set to `Debug`. The root tracer will be set to `Debug`, too.
func QuickConfig(t *testing.T, selectors ...string) func() {
	tracing.RegisterTraceAdapter("test", GetAdapter(t), true)
	c := testconfig.Conf{
		"tracing.adapter": "test",
		"tracelevel.root": "Debug",
	}
	for _, sel := range selectors {
		c["tracelevel."+sel] = "Debug"
	}
	if err := trace2go.ConfigureRoot(c, "tracelevel", trace2go.ReplaceTracers(true)); err != nil {
		t.Fatal(err)
	}
	tracing.SetTraceSelector(trace2go.Selector())
	return trace2go.Teardown
}
