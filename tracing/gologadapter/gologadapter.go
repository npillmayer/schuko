/*
Package gologadapter implements tracing with the default Go logger.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
gologadapter uses the default Go logging mechanism.

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>
*/
package gologadapter

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/npillmayer/schuko/tracing"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a Go standard logger.
type Tracer struct {
	log   *log.Logger
	level tracing.TraceLevel
}

var logLevelPrefix = []string{"ERROR ", "INFO  ", "DEBUG "}

// New creates a new Tracer instance based on a Go logger.
func New() tracing.Trace {
	return &Tracer{
		log:   log.New(os.Stderr, logLevelPrefix[0], log.Ltime),
		level: tracing.LevelError,
	}
}

// GetAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter() tracing.Adapter {
	return New
}

// ----------------------------------------------------------------------------

// P is part of interface Trace
func (t *Tracer) P(key string, val any) tracing.Trace {
	l := &logentry{tracer: t}
	l.p = fmt.Sprintf("[%s=%v] ", key, val)
	return l
}

// Debugf is part of interface Trace
func (t *Tracer) Debugf(s string, args ...any) {
	if t.level < tracing.LevelDebug {
		return
	}
	t.output(tracing.LevelDebug, "", s, args...)
}

// Infof is part of interface Trace
func (t *Tracer) Infof(s string, args ...any) {
	if t.level < tracing.LevelInfo {
		return
	}
	t.output(tracing.LevelInfo, "", s, args...)
}

// Errorf is part of interface Trace
func (t *Tracer) Errorf(s string, args ...any) {
	if t.level < tracing.LevelError {
		return
	}
	t.output(tracing.LevelError, "", s, args...)
}

// SetTraceLevel is part of interface Trace
func (t *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	t.level = l
	t.log.SetPrefix(logLevelPrefix[int(l)])
}

// GetTraceLevel is part of interface Trace
func (t *Tracer) GetTraceLevel() tracing.TraceLevel {
	return t.level
}

// SetOutput is part of interface Trace
func (t *Tracer) SetOutput(writer io.Writer) {
	t.log.SetOutput(writer)
}

func (t *Tracer) output(l tracing.TraceLevel, p string, s string, args ...any) {
	t.log.SetPrefix(logLevelPrefix[int(l)])
	if p == "" { // if no prefix present
		t.log.Printf(s, args...)
	} else {
		t.log.Println(fmt.Sprintf(p+s, args...))
	}
}

// ----------------------------------------------------------------------------

// logentry is a helper for prefixed tracing
type logentry struct { // will have to implement tracing.Trace
	tracer *Tracer // tracer where this logentry will go
	p      string  // prefix
}

func (l *logentry) Debugf(s string, args ...any) {
	if l.tracer.level < tracing.LevelDebug {
		return
	}
	l.tracer.output(tracing.LevelDebug, "", s, args...)
}

func (l *logentry) Infof(s string, args ...any) {
	if l.tracer.level < tracing.LevelInfo {
		return
	}
	l.tracer.output(tracing.LevelDebug, "", s, args...)
}

func (l *logentry) Errorf(s string, args ...any) {
	if l.tracer.level < tracing.LevelError {
		return
	}
	l.tracer.output(tracing.LevelDebug, l.p, s, args...)
}

func (l *logentry) P(key string, val any) tracing.Trace {
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
	l.p = l.p + p
	return l
}

func (l *logentry) SetTraceLevel(tracing.TraceLevel)  {}
func (l *logentry) GetTraceLevel() tracing.TraceLevel { return l.tracer.GetTraceLevel() }
func (l *logentry) SetOutput(writer io.Writer)        {}
