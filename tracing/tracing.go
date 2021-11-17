/*
Package tracing decouples modules and applications from concrete logging implementations.
It abstracts away details of application logging.

Logging/tracing/tracing is a cross-cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling: higher-level
classes/packages are infected with log classes/packages.
That is relevant especially in the context of main applications depending
on external supporting modules, where these modules might want to perform
logging in a specific way incompatible with the main application. For example,
the main application might want to use `logrus` logging, while a supporting
external module logs using the Go standard logger. This is the reason for
the existence of `commons-logging` and the `log4j2` API-definition in the Java
world.

Adapters to concrete logging-implementations are made available by sub-packages of
package `tracing`. The core package `tracing` does not bind to any specific
implementation and creates no dependencies. Deciding for a concrete logger/tracer
is completely up to the main application, where it's perfectly okay to create
a logger/tracer-dependency.

Resources

https://dave.cheney.net/2015/11/05/lets-talk-about-logging

https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern


License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package tracing

import (
	"errors"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/npillmayer/schuko"
)

// TraceLevel is a type for leveled tracing.
// All concrete Tracer implementations will support trace-levels.
type TraceLevel uint8

// We support three trace levels.
const (
	LevelError TraceLevel = iota
	LevelInfo
	LevelDebug
)

func (tl TraceLevel) String() string {
	switch tl {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelError:
		return "Error"
	}
	return "<unknown>"
}

// TraceLevelFromString will find a trace level from a string.
// It will recognize "Debug", "Info" and "Error". Default is
// LevelInfo, if `sl` is not recognized.
//
// String comparison is case-insensitive.
func TraceLevelFromString(sl string) TraceLevel {
	switch strings.ToLower(sl) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "error":
		return LevelError
	}
	return LevelInfo // default
}

// Trace is an interface to be implemented by a concrete tracing adapter.
// For examples please refer to the sub-packages of package tracing.
//
// Tracers should support parameter/field tracing given by P(...).
// An example would be
//
//    tracer.P("mycontext", "value").Debugf("message within context")
//
// Tracers should be prepared to trace to console as well as to a file.
// By convention, no newlines at the end of tracing messages will be passed
// by clients.
type Trace interface {
	Errorf(string, ...interface{}) // trace on error level
	Infof(string, ...interface{})  // trace on level ≥ info
	Debugf(string, ...interface{}) // trace on level ≥ debug
	P(string, interface{}) Trace   // parameter/context tracing
	SetTraceLevel(TraceLevel)      // change the trace level
	GetTraceLevel() TraceLevel     // get the currently active trace level
	SetOutput(io.Writer)           // route tracing output to a writer
}

// Tracefile is the global file where tracing goes to.
// If tracing goes to a file (globally), variable Tracefile should
// point to it. It need not be set if tracing goes to console.
//
// Deprecated: Please use your own app-wide variable.
var Tracefile *os.File // deprecated, will be removed with V1

type TraceSelector interface {
	Select(which string) Trace
}

// SetTraceFactory sets a global TraceSelector.
//
// The use of a global TraceSelector is not mandatory. The default implementation
// returns a no-op tracer for every call.
//
// See also function Select.
func SetTraceSelector(sel TraceSelector) {
	selectorMutex.Lock()
	defer selectorMutex.Unlock()
	selector = sel
}

var selector TraceSelector
var selectorMutex = &sync.RWMutex{} // guard selector

type selectnoOpTracer struct{}

func (snop selectnoOpTracer) Select(string) Trace {
	return noOpTrace{}
}

// Select returns a Trace instance for a given key.
// Initially a default implementation of a TraceSelector is installed which will
// return a no-op tracer for every call, even for key "root".
//
// The use of a global TraceSelector is not mandatory.
func Select(key string) Trace {
	selectorMutex.RLock()
	defer selectorMutex.RUnlock()
	if selector != nil {
		return selector.Select(key)
	}
	return selectnoOpTracer{}.Select(key)
}

// Adapter is a factory function to create a Trace instance.
type Adapter func() Trace

var knownTraceAdapters = map[string]Adapter{
	"nop": func() Trace {
		return noOpTrace{}
	},
}
var adapterMutex = &sync.RWMutex{} // guard knownTraceAdapters[]

// RegisterTraceAdapter is an extension point for clients who want to use
// their own tracing adapter implementation.
// `key` will be used at configuration initialization time to identify
// this adapter, e.g. in configuration files.
//
// Clients will have to call this before any call to tracing-initialization,
// otherwise the adapter cannot be found.
func RegisterTraceAdapter(key string, adapter Adapter, replace bool) {
	adapterMutex.Lock()
	defer adapterMutex.Unlock()
	Infof("registering tracing type %q\n", key)
	current, ok := knownTraceAdapters[key]
	if !ok || current == nil || replace {
		knownTraceAdapters[key] = adapter
	}
}

// GetAdapterFromConfiguration gets the concrete tracing implementation adapter
// from the appcation configuration. If optKey is non-empty it is used for
// looking up the adapter type. Otherwise the default config key is used.
// The default configuration key is "tracing.adapter",
// and if that fails "tracing".
//
// The value must be one of the known tracing adapter keys (see RegisterTraceAdapter).
// If the key is not registered, Adapter
// defaults to a minimalistic (bare bones) tracer.
//
func GetAdapterFromConfiguration(conf schuko.Configuration, optKey string) Adapter {
	adapterPackage := conf.GetString("tracing.adapter")
	if adapterPackage == "" {
		adapterPackage = conf.GetString("tracing")
	}
	adapterMutex.RLock()
	defer adapterMutex.RUnlock()
	adapter, ok := knownTraceAdapters[adapterPackage]
	if !ok || adapter == nil {
		Debugf("no adapter found for tracing type %q\n", adapterPackage)
		adapter = knownTraceAdapters["bare"]
	}
	return adapter
}

// Destination opens a tracing destination as an io.Writer. dest may be one of
//
// a) literals "Stdout" or "Stderr"
//
// b) a file URI ("file: //my.log")
//
// More to come.
//
func Destination(dest string) (io.WriteCloser, error) {
	switch strings.ToLower(dest) {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	}
	u, err := url.Parse(dest)
	if err != nil {
		return os.Stderr, err
	}
	if strings.ToLower(u.Scheme) == "file" {
		fname := u.Path
		if fname == "" {
			fname = u.Host
		}
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			return f, nil
		}
		if errors.Is(err, fs.ErrNotExist) {
			return os.Create(fname)
		} else {
			return nil, err
		}
	}
	return os.Stderr, err
}

// --- Dumping values to trace -----------------------------------------------

// With prepares to dump a data structure to a Trace.
// t may not be nil.
//
// Usage:
//
//     tracing.With(mytracer).Dump(anObject)
//
// Dump accepts interface{};
// it uses 'davecgh/go-spew'.
// Dump(…) will not produce output with t having set a level above LevelDebug.
func With(t Trace) _Dumper {
	return _Dumper{&t}
}

// Helper type for dumping of objects.  Created by calls to With().
type _Dumper struct {
	tracer *Trace
}

// Dump dumps an object using a tracer, in level Debug.
//
// d may not be nil.
func (d _Dumper) Dump(name string, obj interface{}) {
	if (*d.tracer).GetTraceLevel() >= LevelDebug {
		str := spew.Sdump(obj)
		(*d.tracer).Debugf(name + " = " + str)
	}
}

// --- Tracing facade --------------------------------------------------------

// Debugf traces at level LevelDebug to the global default tracer.
// This is part of a global tracing facade.
func Debugf(msg string, args ...interface{}) {
	Select("root").Debugf(msg, args...)
}

// Infof traces at level LevelInfo to the global default tracer.
// This is part of a global tracing facade.
func Infof(msg string, args ...interface{}) {
	Select("root").Infof(msg, args...)
}

// Errorf traces at level LevelError to the global default tracer.
// This is part of a global tracing facade.
func Errorf(msg string, args ...interface{}) {
	Select("root").Errorf(msg, args...)
}

// P performs P on the global default tracer (field tracing).
// Field tracing sets a context for a tracing message.
// This is part of a global tracing facade.
func P(k string, v interface{}) Trace {
	r := Select("root")
	return r.P(k, v)
}

// ---------------------------------------------------------------------------

// NoOpTrace returns a void Trace. This is the default for every global tracer.
// Clients will have to use `SetTraceSelector` to change this.
// This tracer will just do nothing.
func NoOpTrace() Trace {
	return noOpTrace{}
}

type noOpTrace struct{}

func (nt noOpTrace) Debugf(string, ...interface{}) {}
func (nt noOpTrace) Infof(string, ...interface{})  {}
func (nt noOpTrace) Errorf(string, ...interface{}) {}
func (nt noOpTrace) SetTraceLevel(TraceLevel)      {}
func (nt noOpTrace) GetTraceLevel() TraceLevel     { return LevelError }
func (nt noOpTrace) SetOutput(io.Writer)           {}
func (nt noOpTrace) P(string, interface{}) Trace   { return nt }
