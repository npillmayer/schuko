/*
Package koanfadapter is for application configuration.

All configuration is started explicitely with a call to

	schuko.Initialize(koanfadapter.New(k)).

There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

*/
package koanfadapter

import (
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/npillmayer/schuko"
)

// KConf respresents a koanf.Koanf configuration.
type KConf struct {
	k        *koanf.Koanf
	tag      string
	suffixes []string
}

// New creates a new koanf configuration adapter.
// If appTag is given, configuration files are searched for at system-dependent
// standard locations, using the appTag for indentification.
//
// Clients can suppress this behaviour by providing an empty appTag.
// See schuko.LocateConfig for details.
//
func New(k *koanf.Koanf, appTag string, suffixes []string) *KConf {
	return &KConf{
		k:        k,
		tag:      appTag,
		suffixes: suffixes,
	}
}

// Init is usually called by schuko.Initialize()
func (c *KConf) Init() {
	c.InitDefaults()
	if c.tag != "" {
		c.InitFromDefaultFile()
	}
}

// InitDefaults is usually called by Init().
func (c *KConf) InitDefaults() {
	c.k.Load(confmap.Provider(map[string]interface{}{
		"tracing": "go",
	}, "."), nil)
}

// InitFromDefaultFile searches at OS-dependent
// “natural” configuration locations for a readable configuration file.
// See schuko.LocateConfig for details.
// Clients can suppress this behaviour by providing an empty appTag or
// an empty suffixes array during creation of the adapter.
//
// InitFromDefaultFile is usually called by Init().
//
func (c *KConf) InitFromDefaultFile() {
	ok, files := schuko.LocateConfig(c.tag, "", c.suffixes)
	if !ok {
		return
	}
	for _, path := range files {
		ext := filepath.Ext(path)
		switch ext {
		case ".jsn", ".json":
			panic(fmt.Sprintf("do not know how to decode %q-files (%q)", ext, path))
		case ".yml", "yaml":
			panic(fmt.Sprintf("do not know how to decode %q-files (%q)", ext, path))
		case ".tml", ".toml":
			panic(fmt.Sprintf("do not know how to decode %q-files (%q)", ext, path))
		case ".nt":
			//TODO
		default:
			panic(fmt.Sprintf("do not know how to decode %q-files (%q)", ext, path))
		}
	}
}

// Set overrides any configuration values set from the environment.
func (c *KConf) Set(key string, value interface{}) {
	c.k.Load(confmap.Provider(map[string]interface{}{
		key: value,
	}, "."), nil)
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c *KConf) IsSet(key string) bool {
	return c.k.Exists(key)
}

// GetString returns a configuration property as a string.
func (c *KConf) GetString(key string) string {
	return c.k.String(key)
}

// GetInt returns a configuration property as an integer.
func (c *KConf) GetInt(key string) int {
	return c.k.Int(key)
}

// GetBool returns a configuration property as a boolean value.
func (c *KConf) GetBool(key string) bool {
	return c.k.Bool(key)
}

// IsInteractive is a predicate: are we running in interactive mode?
func (c *KConf) IsInteractive() bool {
	return true
}

var _ schuko.Configuration = &KConf{}
