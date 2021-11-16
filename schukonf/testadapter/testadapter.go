/*
Package testadapter is for application configuration during tests (deprecated).

Objects of this package may be used by clients directly, but most of the time
they will be instantiated transparently by calls to `testconfig.QuickConfig`.
Clients will usually follow a pattern along the line of:

    import "github.com/npillmayer/schuko/testconfig"

    func TestSomething(t *testing.T) {
         teardown := testconfig.QuickConfig(t)
         defer teardown()
         …
     }

There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.


License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package testadapter

import (
	"os"
	"strconv"
	"strings"

	"github.com/npillmayer/schuko"
)

// Conf represents a lightweight configuration suited for testing.
//
// Deprecated: Use testconfig Conf instead.
//
type Conf struct {
	values map[string]string
}

// New creates a new configuration suited for testing.
//
// Deprecated: Use testconfig New instead.
//
func New() *Conf {
	return &Conf{values: make(map[string]string)}
}

// Initialize initializes a configuration, populating it with defaullt values.
// It just calls `InitDefaults`.
func (c *Conf) Initialize() {
	c.InitDefaults()
}

// InitDefaults is called to fill the test-configuration with sensible defaults
// for testing. It will set the default tracer to a gotestingadapter instance.
func (c *Conf) InitDefaults() {
	m := c.values
	m["tracing.adapter"] = "test"
	etc := os.Getenv("GOPATH") + "/etc"
	m["etc-dir"] = etc
}

// Set overrides the config value for key.
func (c *Conf) Set(key string, value string) (oldval string) {
	oldval = c.values[key]
	c.values[key] = value
	return
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c *Conf) IsSet(key string) bool {
	_, found := c.values[key]
	return found
}

// GetString is part of the interface Configuration
func (c *Conf) GetString(key string) string {
	v := c.values[key]
	return v
}

// GetInt is part of the interface Configuration
func (c *Conf) GetInt(key string) int {
	v, found := c.values[key]
	if !found {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

// GetBool is part of the interface Configuration
func (c *Conf) GetBool(key string) bool {
	v, found := c.values[key]
	if !found {
		return false
	}
	return strings.EqualFold(v, "true")
}

// IsInteractive is a predicate: are we running in interactive mode?
//
// Deprecated: A custom configuration key should be used instead.
func (c *Conf) IsInteractive() bool { return false }

var _ schuko.Configuration = &Conf{}
