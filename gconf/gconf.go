/*
Package gconf initializes a global application configuration (deprecated).

Configuration

All configuration is started explicitely with a call to Initialize().
There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

Tracing

During configuration all global tracers are set up.
To use a concrete logging implementation,
clients will have to use/implement an adapter to tracing.Trace (please
refer to the documentation for package tracing as well as to implementations
of adapters, e.g. for Go log and for logrus).


License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package gconf

import (
	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/gtrace"
	"github.com/npillmayer/schuko/tracing"
)

// The global configuration is initially set up to be a no-op.
var globalConf schuko.Configuration = noconfig{}

type noconfig struct{}

func (nc noconfig) InitDefaults()               {}
func (nc noconfig) IsSet(key string) bool       { return false }
func (nc noconfig) GetString(key string) string { return "" }
func (nc noconfig) GetInt(key string) int       { return 0 }
func (nc noconfig) GetBool(key string) bool     { return false }
func (nc noconfig) IsInteractive() bool         { return false }

// --- Wire Golabl Configuration ----------------------------------------------------

// Initialize is the top level function for setting up the
// application configuration.
// It will call InitDefaults() on the Configuration passed as an argument, and
// make the Configuration available globally.
// Functions in this package will serve as a facade to the Configuration.
func Initialize(conf schuko.Configuration) {
	globalConf = conf
	globalConf.InitDefaults()
	InitTracing(tracing.GetAdapterFromConfiguration(conf))
}

// InitTracing sets up all the global module tracers, reading trace levels
// and tracing destinations from the application configuration.
//
// InitTracing is usually not called directly, but rather called by Initialize().
func InitTracing(adapter tracing.Adapter) {
	gtrace.CreateTracers(adapter)
	ConfigureTracing("")
}

// ConfigureTracing sets up the global tracers using default configuration values.
//
// It is exported as it may be useful in testing scenarios.
func ConfigureTracing(inputfilename string) {
	SetDefaultTracingLevels() // set default trace levels from configuration
	// if GetBool("tracingonline") {
	// 	if inputfilename != "" {
	// 		// do nothing any more
	// 	}
	// }
	gtrace.InterpreterTracer.P("level", gtrace.InterpreterTracer.GetTraceLevel()).Infof("Interpreter-Trace is alive")
	gtrace.CommandTracer.P("level", gtrace.CommandTracer.GetTraceLevel()).Infof("Command-Trace is alive")
	gtrace.EquationsTracer.P("level", gtrace.EquationsTracer.GetTraceLevel()).Infof("Equations-Trace is alive")
	gtrace.SyntaxTracer.P("level", gtrace.SyntaxTracer.GetTraceLevel()).Infof("Syntax-Trace is alive")
	gtrace.GraphicsTracer.P("level", gtrace.GraphicsTracer.GetTraceLevel()).Infof("Graphics-Trace is alive")
	gtrace.ScriptingTracer.P("level", gtrace.ScriptingTracer.GetTraceLevel()).Infof("Scripting-Trace is alive")
	gtrace.CoreTracer.P("level", gtrace.CoreTracer.GetTraceLevel()).Infof("Core-Trace is alive")
	gtrace.EngineTracer.P("level", gtrace.CoreTracer.GetTraceLevel()).Infof("Engine-Trace is alive")
}

// SetDefaultTracingLevels sets all global tracers to their default trace levels,
// read from the application configuration.
func SetDefaultTracingLevels() {
	gtrace.InterpreterTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinginterpreter")))
	gtrace.CommandTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcommands")))
	gtrace.EquationsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingequations")))
	gtrace.SyntaxTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingsyntax")))
	gtrace.GraphicsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinggraphics")))
	gtrace.ScriptingTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingscripting")))
	gtrace.CoreTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcore")))
	gtrace.EngineTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingengine")))
}

// --- Gobal Configuration Facade ---------------------------------------------------

// IsSet is a predicate wether a global configuration property is set.
func IsSet(key string) bool {
	return globalConf.IsSet(key)
}

// GetString returns a global configuration property as a string.
func GetString(key string) string {
	return globalConf.GetString(key)
}

// GetInt returns a global configuration property as an integer.
func GetInt(key string) int {
	return globalConf.GetInt(key)
}

// GetBool returns a global configuration property as a boolean value.
func GetBool(key string) bool {
	return globalConf.GetBool(key)
}

// IsInteractive is a predicate: are we running in interactive mode?
func IsInteractive() bool {
	return globalConf.IsInteractive()
}
