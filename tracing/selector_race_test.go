package tracing

import (
	"runtime"
	"sync"
	"testing"
)

type selectorRaceTrace struct {
	noOpTrace
	id int
}

func TestSelectConcurrentReplacementAndTeardown(t *testing.T) {
	// This test owns the global selector until all workers have stopped.
	selectorMutex.RLock()
	previous := selector
	selectorMutex.RUnlock()
	defer SetTraceSelector(previous)
	SetTraceSelector(nil)
	if got := Select("root"); got != globalNoOpTrace {
		t.Fatal("nil selector must return no-op tracer")
	}
	a, b := &selectorRaceTrace{id: 1}, &selectorRaceTrace{id: 2}
	selectors := []TraceSelector{genericSelector{tracer: a}, nil, genericSelector{tracer: b}}
	for _, sel := range selectors {
		SetTraceSelector(sel)
		want := Trace(globalNoOpTrace)
		if sel != nil {
			want = sel.Select("child")
		}
		if got := Select("child"); got != want {
			t.Fatal("installed selector was not used")
		}
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 2000; i++ {
			SetTraceSelector(selectors[i%len(selectors)])
			runtime.Gosched()
		}
	}()
	for reader := 0; reader < 8; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < 4000; i++ {
				got := Select("worker")
				if got != a && got != b && got != globalNoOpTrace {
					t.Errorf("unexpected tracer %T", got)
					return
				}
				got.Debugf("safe during selector replacement")
			}
		}()
	}
	close(start)
	wg.Wait()
	SetTraceSelector(nil)
	if got := Select("worker"); got != globalNoOpTrace {
		t.Fatal("teardown must restore no-op selection")
	}
}
