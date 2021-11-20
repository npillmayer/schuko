package trace2go

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/npillmayer/schuko"
	"github.com/npillmayer/schuko/tracing"
)

// TODO override root.SetOutput to redirect all children ?
// Better: let all children be configured thru the same config

// mx guards access to the root tracer
var mx sync.RWMutex

// initRoot guards auto-init of the root tracer
var initRoot sync.Once

// we manage a global root tracer
var root tracing.Trace

// Root returns a reference to the application-global root tracer.
//
// If the root tracer has not yet been initialized, this function will do so before
// returning a reference to it. However, clients will usually want to initialize the
// root tracer in a controlled fashion by calling `ConfigureRoot(…)` first.
//
// If the root tracer has to be auto-initialized (i.e., without a prior call to
// InitRoot), it will default to an opaque tracing.Trace implementation with rudimentary
// functionality:
//
// a) it only traces errors (LevelError)
//
// b) it prints to os.Stderr by default, but reacts to SetOutput(…)
//
// c) no field tracing is available
//
// Later calls to ConfigureRoot can substitute the root tracer.
//
func Root() tracing.Trace {
	mx.RLock()
	defer mx.RUnlock()
	if root == nil {
		initRoot.Do(func() {
			root = &_BareBonesTrace{}
		})
	}
	return root
}

// ConfigureRoot configures the root tracer, given configuration conf.
func ConfigureRoot(conf schuko.Configuration, prefixKey string, opts ...RootOption) error {
	var err error
	r := newRootTracer(conf, prefixKey)
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return err
		}
	}
	r.init()
	mx.Lock()
	defer mx.Unlock()
	_, isBB := root.(*_BareBonesTrace)
	if root == nil || isBB {
		root = r
		//root.Infof("welcome to the new root tracer")
	} else {
		root.Infof("replacing root tracer")
		root = r
		root.Infof("welcome to the new root tracer")
		if r.replaceChildren {
			r := root.(*rootTracer)
			childMx.Lock()
			defer childMx.Unlock()
			// 'range selectableTracers' creates a possible race condition
			// therefore we lock childMx
			for k := range selectableTracers {
				if k == "root" {
					continue
				}
				r.Errorf("replacing tracer \"%s\"", k)
				ch := r.adapter()
				prevCh := setTracer(k, ch)
				prevCh.Infof("replacing this tracer")
				ch.Infof("welcome to the new tracer")
			}
		}
	}
	return err
}

// RootOption is a type to influence initialization of the root tracer.
// Multiple options may be passed to `ConfigureRoot(…)`.
type RootOption _RootOption

type _RootOption func(*rootTracer) error

// ReplaceTracers, when set, will stop all active tracers and replace them
// with new ones which are configured from the new root tracer.
//
// Use it like this:
//
//    err := ConfigureRoot(myconf, "", ReplaceTracers(true))
//
// New tracers replacing existing ones will inherit their trace level.
//
func ReplaceTracers(replace bool) RootOption {
	return func(r *rootTracer) error {
		r.replaceChildren = replace
		return nil
	}
}

// AdapterKey will set a configuration key which, during initialization, will be used
// to search for an adapter type. The key may optionally be set within the configuration passed
// as an argument to ConfigureRoot.
//
// Use it like this:
//
//    err := ConfigureRoot(myconf, "", AdapterKey("look.for.this.adapter"))
//
// with a configuration setting of:
//
//    look.for.this.adapter: logrus
//
// Please see also tracing.GetAdapterFromConfiguration
//
func AdapterKey(key string) RootOption {
	return func(r *rootTracer) error {
		r.optAdapterKey = key
		return nil
	}
}

// --- Root tracer type ------------------------------------------------------

type rootTracer struct {
	tracing.Trace
	config          schuko.Configuration
	prefixKey       string
	optAdapterKey   string
	adapter         tracing.Adapter
	replaceChildren bool
}

// newRootTracer creates a rootTracer struct and populates it with values
// found in configuration `conf`. A prefixKey – if present – is prepended to
// search for configuration values. For example, having these config keys
// loaded:
//
//    trace.root:            Info
//    trace.myapp.mymodule:  Debug
//
// prefixKey should be set to "trace".
//
func newRootTracer(conf schuko.Configuration, prefixKey string) *rootTracer {
	t := &rootTracer{
		config:    conf,
		prefixKey: prefixKey,
	}
	return t
}

func (t *rootTracer) init() {
	adapter := tracing.GetAdapterFromConfiguration(t.config, t.optAdapterKey)
	if adapter == nil {
		adapter = func() tracing.Trace {
			return &_BareBonesTrace{}
		}
	}
	t.adapter = adapter // remember it for child traces
	t.Trace = adapter()
	if l := getValue(t.config, t.prefixKey, "root"); l != "" {
		t.SetTraceLevel(tracing.TraceLevelFromString(l))
	}
}

// --- Integrate as tracing.Selector -----------------------------------------

func Selector() tracing.TraceSelector {
	return selector(trace2goSelector)
}

func trace2goSelector(name string) tracing.Trace {
	return GetOrCreateTracer(name)
}

type selector func(string) tracing.Trace

func (sel selector) Select(name string) tracing.Trace {
	return sel(name)
}

var _ tracing.TraceSelector = selector(trace2goSelector)

// --- Children tracers ------------------------------------------------------

// We will manage a map of keys -> tracers.
var childMx *sync.RWMutex = &sync.RWMutex{}
var selectableTracers map[string]tracing.Trace = make(map[string]tracing.Trace, 10)

// GetTracer returns the tracer associated with name, if any.
func GetTracer(name string) tracing.Trace {
	childMx.RLock()
	defer childMx.RUnlock()
	if child, ok := selectableTracers[name]; ok {
		return child
	}
	return nil
}

// GetOrCreateTracer returns the tracer associated with name. If none is set,
// a new tracer is created and associated with `name`.
func GetOrCreateTracer(name string) tracing.Trace {
	t := GetTracer(name)
	if t == nil {
		t, _ = NewTracer(name, false)
	}
	return t
}

// NewTracer associates a new tracer with a name. Returns the tracer
// occupying the slot, if any.
//
// If parameter `replace` is true, a new tracer will replace an existing one
// for this name.
//
func NewTracer(name string, replace bool) (tracing.Trace, tracing.Trace) {
	var trace tracing.Trace
	if r, ok := Root().(*rootTracer); ok {
		trace = r.adapter()
		level := getValue(r.config, r.prefixKey, name)
		trace.SetTraceLevel(tracing.TraceLevelFromString(level))
	} else {
		return Root(), nil
	}
	childMx.Lock()
	defer childMx.Unlock()
	if prev, ok := selectableTracers[name]; ok { // if exits one...
		if !replace {
			return trace, prev
		}
		selectableTracers[name] = trace // ...replace it
		return trace, prev
	}
	selectableTracers[name] = trace
	return trace, nil
}

// setTracer associates a tracer with a name. Returns the tracer previously
// occupying the slot, if any.
//
// New tracers replacing existing ones will inherit their trace level.
//
// Not protected by childMx.
//
func setTracer(name string, trace tracing.Trace) tracing.Trace {
	prev, ok := selectableTracers[name]
	selectableTracers[name] = trace
	if !ok {
		return prev
	}
	trace.SetTraceLevel(prev.GetTraceLevel())
	return prev
}

// Teardown removes the trace2go root tracer and any existing child tracers,
// and detaches trace2go from the tracing-facade (`tracing.Select(…)`).
func Teardown() {
	mx.Lock()
	defer mx.Unlock()
	fmt.Printf("Tearing down trace2go tracing\n")
	childMx.Lock()
	defer childMx.Unlock()
	selectableTracers = make(map[string]tracing.Trace)
	root = nil
	tracing.SetTraceSelector(nil)
}

// --- Bare bones tracer -----------------------------------------------------

// _BareBonesTrace is a minimalistic Trace implementation. Its use is mainly
// for stepping in in case no tracer is intialized yet.
// This may be especially the case during application-initialization.
//
type _BareBonesTrace struct {
	out io.Writer
}

// Tracer implements tracing.Selector
func (bbt *_BareBonesTrace) Tracer(name string) tracing.Trace {
	if name == "" || name == "root" {
		return bbt
	}
	return tracing.NoOpTrace()
}

// Debugf does nothing
func (bbt _BareBonesTrace) Debugf(string, ...interface{}) {}

// Infof does nothing
func (bbt _BareBonesTrace) Infof(string, ...interface{}) {}

// Errorf traces errors. It is the only trace level implemented for
// bare bones tracers.
func (bbt _BareBonesTrace) Errorf(msg string, args ...interface{}) {
	if bbt.out == nil {
		fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
		return
	}
	fmt.Fprintf(bbt.out, msg, args...)
}

// P does nothing
func (bbt *_BareBonesTrace) P(string, interface{}) tracing.Trace { return bbt }

// SetTraceLevel does nothing (trace level is always LevelError)
func (bbt _BareBonesTrace) SetTraceLevel(tracing.TraceLevel) {}

// GetTraceLevel returns LevelError.
func (bbt _BareBonesTrace) GetTraceLevel() tracing.TraceLevel { return tracing.LevelError }

// SetOutput redirect the output to w. Initially output will be to
// `os.Stderr`.
func (bbt *_BareBonesTrace) SetOutput(w io.Writer) {
	bbt.out = w
}

// ---------------------------------------------------------------------------

func getValue(conf schuko.Configuration, prefixKey string, key string) string {
	k := prefixKey + key
	if v := conf.GetString(k); v != "" {
		return v
	}
	k = prefixKey + "." + key
	if v := conf.GetString(k); v != "" {
		return v
	}
	k = prefixKey + "/" + key
	if v := conf.GetString(k); v != "" {
		return v
	}
	return ""
}
