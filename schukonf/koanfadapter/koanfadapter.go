/*
Package koanfadapter is for application configuration with knadh/koanf.

All configuration is started explicitely with a call to

	conf := koanfadapter.New(…)
	conf.InitDefaults(conf)

There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © 2017–2021 Norbert Pillmayer <norbert@pillmayer.com>

# Parser

NestedTextParser is a thin wrapper on top of npillmayer/nestext to enable using
NestedText (see https://nestedtext.org) as a configuration format.
Koanf needs parsers to implement the koanf.Parser interface.

The parser is also available at `knadh/koanf.parsers.nestedtext`.
*/
package koanfadapter

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/npillmayer/schuko"
)

// KConf respresents a koanf.Koanf configuration.
type KConf struct {
	k        *koanf.Koanf
	tag      string
	suffixes []string
}

// New creates a new koanf configuration adapter. If k is nil, a new Koanf
// will be created, with '.' as an item separator.
//
// If appTag and suffixes are given, configuration files are searched for at
// system-dependent standard locations, using the appTag for indentification
// (please refer to the description of koanfadapter.InitDefaults).
func New(k *koanf.Koanf, appTag string, suffixes []string) *KConf {
	if k == nil {
		k = koanf.New(".")
	}
	return &KConf{
		k:        k,
		tag:      appTag,
		suffixes: suffixes,
	}
}

// Koanf returns the embedded Koanf configuration-object.
func (c *KConf) Koanf() *koanf.Koanf {
	return c.k
}

// Init is usually called by schuko.Initialize()
// func (c *KConf) Init() {
//     c.InitDefaults()
// }

// InitDefaults loads initial configuration settings. For koanf, it does two things:
//
// (1) It sets the default tracer to the builtin Go log
//
// (2) It loads application-specific configuration from
//
//	“natural” configuration locations, if any are found. It does this by calling
//	`InitFromDefaultFile()`
func (c *KConf) InitDefaults() {
	c.k.Load(confmap.Provider(map[string]any{
		"tracing.adapter": "go",
	}, c.k.Delim()), nil)
	if c.tag != "" {
		c.InitFromDefaultFile()
	}
}

// InitFromDefaultFile searches at OS-dependent
// “natural” configuration locations for a readable configuration file.
// See schuko.LocateConfig for details.
// Clients can suppress this behaviour by providing an empty appTag or
// an empty suffixes array during creation of the adapter.
//
// InitFromDefaultFile is usually not called directly by clients, but
// rather by InitDefaults. It is made public to enable clients to override
// it.
func (c *KConf) InitFromDefaultFile() {
	files := schuko.LocateConfig(c.tag, "", c.suffixes)
	if len(files) == 0 {
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
			//TODO remove dependency to go log
			if err := c.k.Load(file.Provider(path), Parser()); err != nil {
				log.Fatalf("error loading NestedText format: %v", err)
			}
		default:
			panic(fmt.Sprintf("do not know how to decode %q-files (%q)", ext, path))
		}
	}
}

// Set overrides any configuration values set from the environment.
func (c *KConf) Set(key string, value any) {
	c.k.Load(confmap.Provider(map[string]any{
		key: value,
	}, c.k.Delim()), nil)
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
//
// Deprecated: A custom configuration key should be used instead.
func (c *KConf) IsInteractive() bool {
	return true
}

var _ schuko.Configuration = &KConf{}
