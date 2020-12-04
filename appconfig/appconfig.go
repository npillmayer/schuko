/*
Package appconfig provides configuration and tracing suitable for
full fledged applications.

BSD License

Copyright (c) 2017–21, Norbert Pillmayer

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package appconfig

import (
	"fmt"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/gconf"
	"github.com/npillmayer/schuko/tracing"
	"github.com/npillmayer/schuko/tracing/logrusadapter"
	"github.com/npillmayer/schuko/viperadapter"
)

// QuickConfig sets up a configuration suitable for applications.
// It will install a configuration adapter for "github.com/spf13/viper".
// If tracing is not selected from viper (which may have its own default for tracing),
// it defaults to "github.com/Sirupsen/logrus".
func QuickConfig() {
	conf := viperadapter.New()
	if !conf.IsSet("tracing") {
		conf.Set("tracing", "logrus")
		adapter := schuko.GetAdapterFromConfiguration(conf)
		if adapter == nil {
			schuko.AddTraceAdapter("logrus", logrusadapter.GetAdapter())
		}
	}
	gconf.Initialize(conf)
}

// WithTracing lets clients select a tracing engine other than the
// default tracing. It will override any tracing selected from the application
// configuration.
func WithTracing(tracekey string, traceadapter tracing.Adapter) (err error) {
	conf := viperadapter.New()
	conf.Set("tracing", tracekey)
	if traceadapter == nil {
		adapter := schuko.GetAdapterFromConfiguration(conf)
		if adapter == nil {
			err = fmt.Errorf("unable to find tracer for key='%s'", tracekey)
		}
	} else {
		schuko.AddTraceAdapter(tracekey, traceadapter)
	}
	gconf.Initialize(conf)
	return
}
