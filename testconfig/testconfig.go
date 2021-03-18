/*
Package testconfig provides configuration and tracing suitable for tests.

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
package testconfig

import (
	"testing"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/gconf"
	"github.com/npillmayer/schuko/testadapter"
	"github.com/npillmayer/schuko/tracing/gotestingadapter"
)

// QuickConfig sets up a configuration suitable for test cases, including tracing.
// It returns a teardown function which should be called at the end of a test.
// The usual pattern will look like this:
//
//     func TestSomething(t *testing.T) {
//          teardown := testconfig.QuickConfig(t)
//          defer teardown()
//          ...
//      }
//
func QuickConfig(t *testing.T, maps ...map[string]string) func() {
	schuko.AddTraceAdapter("test", gotestingadapter.GetAdapter())
	c := testadapter.New()
	c.Set("tracing", "test")
	for _, m := range maps {
		for k, v := range m {
			c.Set(k, v)
		}
	}
	gconf.Initialize(c)
	return gotestingadapter.RedirectTracing(t)
}
