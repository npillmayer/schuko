/*
Package testconfig provides configuration and tracing suitable for tests
(deprecated).

The usual usage-pattern will look like this:

    func TestSomething(t *testing.T) {
        teardown := testconfig.QuickConfig(t)
        defer teardown()
        …
     }

License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package testconfig

import (
	"testing"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/gconf"
	"github.com/npillmayer/schuko/schukonf/testadapter"
	"github.com/npillmayer/schuko/tracing/gotestingadapter"
)

// QuickConfig sets up a configuration suitable for test cases, including tracing.
// It returns a teardown function which should be called at the end of a test.
// The usual pattern will look like this:
//
//     func TestSomething(t *testing.T) {
//          teardown := testconfig.QuickConfig(t, map[string]string {
//              "my-key": "my override value just for testing",
//          })
//          defer teardown()
//          …
//      }
//
// Deprecated: This moves to schukonf testconfig.
//
func QuickConfig(t *testing.T, maps ...map[string]string) func() {
	schuko.AddTraceAdapter("test", gotestingadapter.GetAdapter(t))
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
