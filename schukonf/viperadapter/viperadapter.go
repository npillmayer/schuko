/*
Package viperadapter is for application configuration with spf13/viper.

All configuration is started explicitely with a call to

	conf := viperadapter.New("myapp")
	conf.InitDefaults()

There is no init-call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>
*/
package viperadapter

import (
	"fmt"
	"os"

	"github.com/npillmayer/schuko"
	"github.com/spf13/viper"
)

// VConf respresents a Viper configuration.
type VConf struct {
	name string
}

// New creates a new Viper configuration adapter.
// `name` is used as a tag to locate application-configuration.
func New(name string) *VConf {
	return &VConf{name: name}
}

// Init is usually called by schuko.Initialize()
func (c *VConf) Init() {
	c.InitDefaults()
	c.InitConfigPath()
}

// InitDefaults sets up
func (c *VConf) InitDefaults() {
	viper.SetDefault("tracing", "go")
	viper.SetDefault("tracingonline", true)
}

// InitConfigPath is usually called by Init().
// It sets up a standard config path and searches for application configuration
// files.
func (c *VConf) InitConfigPath() {
	viper.SetConfigName(c.name)                    // name of config file (without extension)
	viper.AddConfigPath(".")                       // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/." + c.name)        // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.config/" + c.name) // call multiple times to add many search paths
	err := viper.ReadInConfig()                    // Find and read the config file
	if err != nil {                                // Handle errors reading the config file
		fmt.Fprintf(os.Stderr, "error reading config file: %s", err.Error())
	}
}

// Set overrides any configuration values.
func (c *VConf) Set(key string, value any) {
	viper.Set(key, value)
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c *VConf) IsSet(key string) bool {
	return viper.IsSet(key)
}

// GetString returns a configuration property as a string.
func (c *VConf) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns a configuration property as an integer.
func (c *VConf) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a configuration property as a boolean value.
func (c *VConf) GetBool(key string) bool {
	return viper.GetBool(key)
}

// IsInteractive is a predicate: are we running in interactive mode?
//
// Deprecated: A custom configuration key should be used instead.
func (c *VConf) IsInteractive() bool {
	return viper.GetBool("tracingonline")
}

var _ schuko.Configuration = &VConf{}
