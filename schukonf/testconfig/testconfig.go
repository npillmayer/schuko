/*
Package testconfig provides configuration and tracing suitable for tests.

The usual usage-pattern will look like this:

	func TestSomething(t *testing.T) {
	    teardown := testconfig.QuickConfig(t)
	    defer teardown()
	    …
	 }

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>
*/
package testconfig

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/npillmayer/schuko"
)

// Conf represents a lightweight configuration suited for testing.
type Conf map[string]any

// InitDefaults is called to fill the test-configuration with sensible defaults
// for testing. It will set the default tracer to a gotestingadapter instance.
//
// Does nothing.
//
// Deprecated: InitDefaults will be removed.
func (c Conf) InitDefaults() {
}

// Set overrides the config value for key.
func (c Conf) Set(key string, value string) (oldval string) {
	oldval = fmt.Sprintf("%v", c[key])
	c[key] = value
	return
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c Conf) IsSet(key string) bool {
	_, found := c[key]
	return found
}

// GetString is part of the interface Configuration
func (c Conf) GetString(key string) string {
	v, found := c[key]
	if !found {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// GetInt is part of the interface Configuration
func (c Conf) GetInt(key string) int {
	v, found := c[key]
	if !found {
		return 0
	}
	var n int
	switch x := v.(type) {
	case string:
		n, _ = strconv.Atoi(x)
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		n, _ = x.(int)
	default:
		n = 0
	}
	return n
}

// GetBool is part of the interface Configuration
func (c Conf) GetBool(key string) bool {
	v, found := c[key]
	if !found {
		return false
	}
	var b bool
	switch x := v.(type) {
	case string:
		b = strings.EqualFold(x, "true")
	case bool:
		b = x
	default:
		b = false
	}
	return b
}

// IsInteractive is a predicate: are we running in interactive mode?
//
// Deprecated: A custom configuration key should be used instead.
func (c Conf) IsInteractive() bool { return false }

var _ schuko.Configuration = &Conf{}
