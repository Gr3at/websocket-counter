package httpsrv

import (
	"testing"
)

func TestStats(t *testing.T) {
	server := New(make(chan string))

	server.incStats("watcher1")
	server.incStats("watcher2")
	server.incStats("watcher2")
	server.incStats("watcher1")
	server.incStats("watcher2")

	for _, ws := range server.sessionStats {
		if ws.id == "watcher1" && ws.sent != 2 {
			t.Fatalf("Expected 2 messages for watcher1, got %d", ws.sent)
		} else if ws.id == "watcher2" && ws.sent != 3 {
			t.Fatalf("Expected 3 messages for watcher2, got %d", ws.sent)
		}
	}

}
