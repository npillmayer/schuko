/*
Package goslogadapter implements tracing with Go's slog package.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
goslogadapter uses the standard "log/slog" mechanism.

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>
*/
package goslogadapter

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/npillmayer/schuko/tracing"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a slog logger.
type Tracer struct {
	log   *slog.Logger
	level *slog.LevelVar
}

// New creates a new Tracer instance based on slog.
func New() tracing.Trace {
	lv := &slog.LevelVar{}
	lv.Set(slog.LevelError)
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lv})
	return &Tracer{
		log:   slog.New(h),
		level: lv,
	}
}

// GetAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter() tracing.Adapter {
	return New
}

// ----------------------------------------------------------------------------

// P is part of interface Trace.
func (t *Tracer) P(key string, val any) tracing.Trace {
	return &logentry{
		tracer: t,
		attrs:  []any{key, val},
	}
}

// Debugf is part of interface Trace.
func (t *Tracer) Debugf(s string, args ...any) {
	t.output(tracing.LevelDebug, nil, s, args...)
}

// Infof is part of interface Trace.
func (t *Tracer) Infof(s string, args ...any) {
	t.output(tracing.LevelInfo, nil, s, args...)
}

// Errorf is part of interface Trace.
func (t *Tracer) Errorf(s string, args ...any) {
	t.output(tracing.LevelError, nil, s, args...)
}

// SetTraceLevel is part of interface Trace.
func (t *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	t.level.Set(translateTraceLevel(l))
}

// GetTraceLevel is part of interface Trace.
func (t *Tracer) GetTraceLevel() tracing.TraceLevel {
	return translateSlogLevel(t.level.Level())
}

// SetOutput is part of interface Trace.
func (t *Tracer) SetOutput(writer io.Writer) {
	h := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: t.level})
	t.log = slog.New(h)
}

func (t *Tracer) output(l tracing.TraceLevel, attrs []any, s string, args ...any) {
	sl := translateTraceLevel(l)
	ctx := context.Background()
	if !t.log.Enabled(ctx, sl) {
		return
	}
	msg := fmt.Sprintf(s, args...)
	if len(attrs) == 0 {
		t.log.Log(ctx, sl, msg)
		return
	}
	t.log.With(attrs...).Log(ctx, sl, msg)
}

func translateSlogLevel(l slog.Level) tracing.TraceLevel {
	switch {
	case l <= slog.LevelDebug:
		return tracing.LevelDebug
	case l < slog.LevelError:
		return tracing.LevelInfo
	default:
		return tracing.LevelError
	}
}

func translateTraceLevel(l tracing.TraceLevel) slog.Level {
	switch l {
	case tracing.LevelDebug:
		return slog.LevelDebug
	case tracing.LevelInfo:
		return slog.LevelInfo
	case tracing.LevelError:
		return slog.LevelError
	}
	return slog.LevelInfo
}

// ----------------------------------------------------------------------------

// logentry is a helper for context tracing.
type logentry struct {
	tracer *Tracer
	attrs  []any
}

func (l *logentry) Debugf(s string, args ...any) {
	l.tracer.output(tracing.LevelDebug, l.attrs, s, args...)
}

func (l *logentry) Infof(s string, args ...any) {
	l.tracer.output(tracing.LevelInfo, l.attrs, s, args...)
}

func (l *logentry) Errorf(s string, args ...any) {
	l.tracer.output(tracing.LevelError, l.attrs, s, args...)
}

func (l *logentry) P(key string, val any) tracing.Trace {
	l.attrs = append(l.attrs, key, val)
	return l
}

func (l *logentry) SetTraceLevel(tracing.TraceLevel)  {}
func (l *logentry) GetTraceLevel() tracing.TraceLevel { return l.tracer.GetTraceLevel() }
func (l *logentry) SetOutput(io.Writer)               {}

