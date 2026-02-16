/*
Package schuko defines types for application configuration and tracing.
Application configuration is addressed by quite a lot of go libraries out there.
We do not intend to re-invent the wheel, but rather place a layer on top of existing
libraries.  In particular, we'll integrate logging/tracing-configuration, making it
easy to re-configure between development and production use.

There is no init-call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

# Attention

As this package nears V1, some re-structuring happenes. Please look out for
`deprecated` tags.

# License

Governed by a 3-Clause BSD license. License file may be found in the root
folder of this module.

Copyright © Norbert Pillmayer <norbert@pillmayer.com>
*/
package schuko

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Configuration is an interface to be implemented by every configuration
// adapter.
type Configuration interface {
	InitDefaults()               // initialize base set of key/value pairs
	IsSet(key string) bool       // is a config key set?
	GetString(key string) string // get config value as string
	GetInt(key string) int       // get config value as integer
	GetBool(key string) bool     // get config value as boolean
}

var knownTraceAdapters = map[string]any{
	//"go": gologadapter.GetAdapter(),
	//"logrus": logrusadapter.GetAdapter(), // now to be set by AddTraceAdapter()
}
var adapterMutex = &sync.RWMutex{} // guard knownTraceAdapters[]

// AddTraceAdapter is an extension point for clients who want to use
// their own tracing adapter implementation.
// key will be used at configuration initialization time to identify
// this adapter, e.g. in configuration files.
//
// Clients will have to call this before any call to tracing-initialization,
// otherwise the adapter cannot be found.
//
// Deprecated: This moves to package tracing RegisterTraceAdapter.
func AddTraceAdapter(key string, adapter any) {
	adapterMutex.Lock()
	defer adapterMutex.Unlock()
	knownTraceAdapters[key] = adapter
}

// GetAdapterFromConfiguration gets the concrete tracing implementation adapter
// from the appcation configuration. The configuration key name is "tracing".
//
// The value must be one of the known tracing adapter keys.
// Default is an adapter for the Go standard log package.
//
// Deprecated: This moves to package tracing.
func GetAdapterFromConfiguration(conf Configuration) any {
	adapterPackage := conf.GetString("tracing")
	adapterMutex.RLock()
	defer adapterMutex.RUnlock()
	adapter := knownTraceAdapters[adapterPackage]
	// if adapter == nil {
	// 	adapter = gologadapter.GetAdapter() // removed this coupling to Go log
	// }
	return adapter
}

// LocateConfig searches configuration files at “natural” configuration locations, which
// are OS-dependent (see os.UserConfigDir).
// The application is given by a tag name, which will be used to search for
// existing directories and files, and an optional pattern.
// Files will have to match one of:
//
//	<pattern>.<suffix>    // if pattern is given
//	<tag>.<suffix>        // if no pattern
//	config.<suffix>       // if no pattern
//	.<tag>.<suffix>       // for $HOME only and no pattern
//
// Allowed file types are given as argument `suffixes`.
//
// Example: An app uses the tag 'myapp'. On a *nix-system the configuration may
// be searched for at
//
//	$HOME/.config/myapp/config.*
//	$HOME/.config/myapp/myapp.*
//	$HOME/.myapp.*
//
// On MacOS it would be searched for in
//
//	$HOME/Library/Application Support/MyApp/
func LocateConfig(appTag string, pattern string, suffixes []string) []string {
	//
	if appTag == "" || len(suffixes) == 0 {
		return nil
	}
	tag := strings.ToLower(appTag)

	var d []fs.DirEntry
	var found bool
	var dir string
	var dirs []string

	homedir, errH := os.UserHomeDir()
	confdir, err := os.UserConfigDir()

	if err == nil && (errH == nil && confdir != homedir) {
		dir = filepath.Join(confdir, appTag)
		if d, err = os.ReadDir(dir); err == nil {
			if found, dirs = dirMatch(dir, d, tag, pattern, suffixes); found {
				return dirs
			}
		}
		dir = filepath.Join(confdir, tag)
		if d, err = os.ReadDir(dir); err == nil {
			if found, dirs = dirMatch(dir, d, tag, pattern, suffixes); found {
				return dirs
			}
		}
	}
	if errH != nil {
		return nil
	}
	dir = filepath.Join(homedir, ".config", tag)
	if d, err = os.ReadDir(dir); err == nil {
		// look for ~/.config/myapp/*
		if found, dirs = dirMatch(dir, d, tag, pattern, suffixes); found {
			return dirs
		}
	}
	if d, err = os.ReadDir(homedir); err == nil {
		// look for ~/.myapp.*
		if found, dirs = dirMatch(homedir, d, tag, pattern, suffixes); found {
			return dirs
		}
	}
	return nil
}

func dirMatch(dir string, d []fs.DirEntry, tag, pattern string, suffixes []string) (bool, []string) {
	m := []string{}
	glob1 := "config.*"
	glob2 := tag + ".*"
	glob3 := "." + tag + ".*"
	for _, e := range d {
		fname := filepath.Base(e.Name())
		if pattern != "" {
			if fm(pattern, fname) {
				ext := strings.TrimLeft(filepath.Ext(fname), ".")
				for _, s := range suffixes {
					if ext == s {
						m = append(m, filepath.Join(dir, fname))
						break
					}
				}
			}
		} else if fm(glob1, fname) || fm(glob2, fname) || fm(glob3, fname) {
			ext := strings.TrimLeft(filepath.Ext(fname), ".")
			for _, s := range suffixes {
				if ext == s {
					m = append(m, filepath.Join(dir, fname))
					break
				}
			}
		}
	}
	if len(m) > 0 {
		return true, m
	}
	return false, nil
}

func fm(pattern, name string) bool {
	ok, _ := filepath.Match(pattern, name)
	return ok
}
