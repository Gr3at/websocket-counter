package httpsrv

import (
	"goapp/internal/pkg/watcher"
	"testing"
	"time"
)

func TestServerWatcher(t *testing.T) {
	var (
		strChan    = make(chan string)
		httpServer = New(strChan)
		watcher    = watcher.New()
	)

	httpServer.addWatcher(watcher)
	if len(httpServer.watchers) != 1 {
		t.Fatalf("expected 1 watcher, got %d", len(httpServer.watchers))
	}

	if _, exists := httpServer.watchers[watcher.GetWatcherId()]; !exists {
		t.Fatalf("expected watcher with ID %s to be added", watcher.GetWatcherId())
	}

	if err := watcher.Start(); err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	httpServer.notifyWatchers("test message")
	select {
	case counter := <-watcher.Recv():
		if counter.Iteration != 1 {
			t.Fatalf("expected counter.Iteration = 1, got %d", counter.Iteration)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("expected to receive message from watcher")
	}

	httpServer.removeWatcher(watcher)
	if len(httpServer.watchers) != 0 {
		t.Fatalf("expected 0 watchers, got %d", len(httpServer.watchers))
	}
}
