package gotestingadapter

import (
	"fmt"
	"sync"
	"testing"

	"github.com/npillmayer/schuko/tracing"
)

func TestConcurrentLoggingAndLevels(t *testing.T) {
	tr := New(t)
	tr.SetTraceLevel(tracing.LevelDebug)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				tr.Debugf("worker %d", worker)
				tr.Infof("iteration %d", i)
				tr.Errorf("test message %d", i)
				tr.SetTraceLevel(tracing.TraceLevel(i % 3))
				_ = tr.GetTraceLevel()
				tr.P("worker", worker).Infof("contextual message")
			}
		}(worker)
	}
	wg.Wait()
}

func TestConcurrentContextsAreIndependent(t *testing.T) {
	tr := New(t).(*Tracer)
	tr.SetTraceLevel(tracing.LevelDebug)
	const workers = 8
	ready := make(chan *Tracer, workers)
	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := tr.P("worker", i).(*Tracer)
			ready <- ctx
			<-release
			if ctx.p != fmt.Sprintf("[worker=%d] ", i) {
				t.Errorf("context leaked: %q for worker %d", ctx.p, i)
			}
			ctx.Infof("message")
		}(i)
	}
	seen := map[*Tracer]bool{}
	for i := 0; i < workers; i++ {
		seen[<-ready] = true
	}
	if len(seen) != workers || tr.p != "" {
		t.Error("P mutated/shared the receiver")
	}
	close(release)
	wg.Wait()
}

func TestLegacyRedirectConcurrentLogging(t *testing.T) {
	tr := &Tracer{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			done := RedirectTracing(t)
			done()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			tr.Errorf("legacy message")
		}
	}()
	wg.Wait()
}
