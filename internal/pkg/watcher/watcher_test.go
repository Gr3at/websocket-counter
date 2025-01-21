package watcher

import (
	"testing"
)

func TestWatcher(t *testing.T) {
	w := New()
	if w.GetWatcherId() == "" {
		t.Fatal("expected watcher ID to have a value")
	}

	err := w.Start()
	defer w.Stop()
	if err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}

	w.Send("test message")
	<-w.Recv()
	w.Send("test message 2")
	counter := <-w.Recv()
	if counter.Iteration != 2 {
		t.Fatalf("expected counter.Iteration = 2, got %d", counter.Iteration)
	}

	w.ResetCounter()
	counter = <-w.Recv()
	if counter.Iteration != 0 {
		t.Fatalf("expected counter.Iteration = 0 after reset, got %d", counter.Iteration)
	}

	w.Send("test message 3")
	counter = <-w.Recv()
	if counter.Iteration != 1 {
		t.Fatalf("checkig if increment works after reset: expected counter.Iteration = 1, got %d", counter.Iteration)
	}
}
