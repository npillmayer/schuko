/*
Package gtrace lets clients set up a set of global tracers (deprecated).

gtrace infects clients with some semantics for structuring tracing/logging on
an application-wide level.
However, the imposed overhead is very slim, i.e. it boils down to having
some global variables declared (but not necessarily used):

	CommandTracer     : Tracing interactive or batch commands from users
	CoreTracer        : Tracing application core
	ScriptingTracer   : Tracing embedded scripting host(s)
	InterpreterTracer : Tracing DSL interpreting
	SyntaxTracer      : Tracing lexing/parsing of DSLs
	GraphicsTracer    : Tracing graphics routines
	EngineTracer      : Tracing web-engine routines
	EquationsTracer   : Tracing arithmetic

Clients are free to selectively use any of these tracers.
They are initially set up to do nothing (no-ops).
Clients will use the configuration package to set it up in a
meaningful way.


License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package gtrace

import (
	"errors"
	"io"

	"github.com/npillmayer/schuko/tracing"
)

// This is the set of standard module tracers for our application.
//
// All tracers are set up to be no-ops, initially.
// This approach poses a little burden on (selective) clients, but is
// useful for de-coupling the various packages and modules from the
// tracing/logging mechanism.
var (
	EquationsTracer   = NoOpTrace
	InterpreterTracer = NoOpTrace
	SyntaxTracer      = NoOpTrace
	CommandTracer     = NoOpTrace
	GraphicsTracer    = NoOpTrace
	ScriptingTracer   = NoOpTrace
	CoreTracer        = NoOpTrace
	EngineTracer      = NoOpTrace
)

// NoOpTrace is a void Trace. Initially, all tracers will be set up to be no-ops.
// Clients will have to configure concrete tracing backends, usually by calling
// application configuration with a tracing adapter.
var NoOpTrace tracing.Trace = nooptrace{}

type nooptrace struct{}

func (nt nooptrace) Debugf(string, ...interface{})       {}
func (nt nooptrace) Infof(string, ...interface{})        {}
func (nt nooptrace) Errorf(string, ...interface{})       {}
func (nt nooptrace) SetTraceLevel(tracing.TraceLevel)    {}
func (nt nooptrace) GetTraceLevel() tracing.TraceLevel   { return tracing.LevelError }
func (nt nooptrace) SetOutput(io.Writer)                 {}
func (nt nooptrace) P(string, interface{}) tracing.Trace { return nt }

// CreateTracers creates all global tracers, given a function to
// create a concrete Trace instance.
func CreateTracers(newTrace tracing.Adapter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to create global tracers")
		}
	}()
	EquationsTracer = newTrace()
	InterpreterTracer = newTrace()
	SyntaxTracer = newTrace()
	CommandTracer = newTrace()
	GraphicsTracer = newTrace()
	ScriptingTracer = newTrace()
	CoreTracer = newTrace()
	EngineTracer = newTrace()
	return
}

// Mute sets all global tracers to LevelError.
func Mute() {
	InterpreterTracer.SetTraceLevel(tracing.LevelError)
	CommandTracer.SetTraceLevel(tracing.LevelError)
	EquationsTracer.SetTraceLevel(tracing.LevelError)
	SyntaxTracer.SetTraceLevel(tracing.LevelError)
	GraphicsTracer.SetTraceLevel(tracing.LevelError)
	ScriptingTracer.SetTraceLevel(tracing.LevelError)
	CoreTracer.SetTraceLevel(tracing.LevelError)
	EngineTracer.SetTraceLevel(tracing.LevelError)
}
