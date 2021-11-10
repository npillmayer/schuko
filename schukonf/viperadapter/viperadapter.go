/*
Package viperadapter is for application configuration.

All configuration is started explicitely with a call to

	schuko.Initialize(viperadapter.New("myconf")).

There is no init-call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package viperadapter

import (
	"fmt"

	"github.com/npillmayer/schuko"
	"github.com/spf13/viper"
)

// VConf respresents a Viper configuration.
type VConf struct {
	name string
}

// New creates a new Viper configuration adapter.
func New(name string) *VConf {
	return &VConf{name: name}
}

// Init is usually called by schuko.Initialize()
func (c *VConf) Init() {
	c.InitDefaults()
	c.InitConfigPath()
}

// InitDefaults is usually called by Init().
func (c *VConf) InitDefaults() {
	viper.SetDefault("tracing", "go")
	viper.SetDefault("tracingonline", true)
	viper.SetDefault("tracingequations", "Error")
	viper.SetDefault("tracingsyntax", "Error")
	viper.SetDefault("tracingcommands", "Error")
	viper.SetDefault("tracinginterpreter", "Error")
	viper.SetDefault("tracinggraphics", "Error")

	viper.SetDefault("tracingcapsules", "Error")
	viper.SetDefault("tracingrestores", "Error")
	viper.SetDefault("tracingchoices", true)
}

// InitConfigPath is usually called by Init().
func (c *VConf) InitConfigPath() {
	viper.SetConfigName(c.name)             // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AddConfigPath("$GOPATH/etc/")     // path to look for the config file in
	viper.AddConfigPath("$HOME/." + c.name) // call multiple times to add many search paths
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		panic(fmt.Errorf("fatal error reading config file: %s", err.Error()))
	}
}

// Set overrides any configuration values set from the environment.
func (c *VConf) Set(key string, value interface{}) {
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
func (c *VConf) IsInteractive() bool {
	return viper.GetBool("tracingonline")
}

var _ schuko.Configuration = &VConf{}
