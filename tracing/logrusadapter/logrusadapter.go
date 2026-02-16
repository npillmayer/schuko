/*
Package logrusadapter implements tracing with the logrus logger.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
logrus uses "github.com/sirupsen/logrus" as the means for tracing.

# BSD License

# Copyright (c) 2017–21, Norbert Pillmayer

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>
*/
package logrusadapter

import (
	"io"

	"github.com/npillmayer/schuko/tracing"
	"github.com/sirupsen/logrus"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a logrus logger.
type Tracer struct {
	log *logrus.Logger
	p   *logrus.Entry
}

// New creates a new Tracer instance based on a logrus logger.
func New() tracing.Trace {
	return &Tracer{logrus.New(), nil}
}

// NewAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter() tracing.Adapter {
	return New
}

// Interface tracing.Trace
func (t *Tracer) P(key string, val any) tracing.Trace {
	t.p = t.log.WithField(key, val)
	return t
}

// Interface tracing.Trace
func (t *Tracer) Debugf(s string, args ...any) {
	if t.p != nil {
		t.p.Debugf(s, args...)
		t.p = nil
	} else {
		t.log.Debugf(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) Infof(s string, args ...any) {
	if t.p != nil {
		t.p.Infof(s, args...)
		t.p = nil
	} else {
		t.log.Infof(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) Errorf(s string, args ...any) {
	if t.p != nil {
		t.p.Errorf(s, args...)
		t.p = nil
	} else {
		t.log.Errorf(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	t.log.SetLevel(translateTraceLevel(l))
}

// Interface tracing.Trace
func (t *Tracer) GetTraceLevel() tracing.TraceLevel {
	return translateLogLevel(t.log.Level)
}

// Interface tracing.Trace
func (t *Tracer) SetOutput(writer io.Writer) {
	t.log.Out = writer
	t.log.Formatter = &logrus.TextFormatter{}
}

func translateLogLevel(l logrus.Level) tracing.TraceLevel {
	switch l {
	case logrus.DebugLevel:
		return tracing.LevelDebug
	case logrus.InfoLevel:
		return tracing.LevelInfo
	case logrus.ErrorLevel:
		return tracing.LevelError
	}
	return tracing.LevelDebug
}

func translateTraceLevel(l tracing.TraceLevel) logrus.Level {
	switch l {
	case tracing.LevelDebug:
		return logrus.DebugLevel
	case tracing.LevelInfo:
		return logrus.InfoLevel
	case tracing.LevelError:
		return logrus.ErrorLevel
	}
	return logrus.DebugLevel
}
