package httpsrv

import (
	"encoding/json"
	"goapp/internal/pkg/watcher"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHealthHandler(t *testing.T) {
	strChan := make(chan string)
	server := New(strChan)

	req := httptest.NewRequest(http.MethodGet, "/goapp/health", nil)
	rec := httptest.NewRecorder()

	server.handlerHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHomeHandler(t *testing.T) {
	strChan := make(chan string)
	server := New(strChan)

	req := httptest.NewRequest(http.MethodGet, "/goapp", nil)
	rec := httptest.NewRecorder()

	server.handlerHome(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestWebSocketHandler(t *testing.T) {
	strChan := make(chan string)
	server := New(strChan)

	testServer := httptest.NewServer(http.HandlerFunc(server.handlerWebSocket))
	defer testServer.Close()

	address, _ := strings.CutPrefix(testServer.URL, "http")
	wsURL := "ws" + address
	headers := http.Header{}
	headers.Set("Origin", "http://localhost:8080")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		t.Fatalf("failed to establish WebSocket connection: %v", err)
	}
	defer ws.Close()

	resetMessage := watcher.CounterReset{}
	msgBytes, _ := json.Marshal(resetMessage)
	if err := ws.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		t.Fatalf("failed to send WebSocket message: %v", err)
	}

	select {
	case <-time.After(200 * time.Millisecond):
	default:
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read WebSocket message: %v", err)
		}

		var counter watcher.Counter
		if err := json.Unmarshal(p, &counter); err != nil {
			t.Fatalf("failed to unmarshal WebSocket response: %v", err)
		}

		if counter.Iteration != 0 {
			t.Fatalf("expected counter.Iteration = 0 after reset, got %d", counter.Iteration)
		}
	}
}
