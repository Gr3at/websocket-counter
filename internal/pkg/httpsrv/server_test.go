package httpsrv

import (
	"encoding/json"
	"fmt"
	"goapp/internal/pkg/watcher"
	"net/http"
	"net/http/httptest"
	"os"
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
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:8080")
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

type csrfTestCase struct {
	allowedOrigins     string
	origin             string
	expectedStatusCode int
}

func TestWebSocketHandlerCSRF(t *testing.T) {
	strChan := make(chan string)
	server := New(strChan)

	testServer := httptest.NewServer(http.HandlerFunc(server.handlerWebSocket))
	defer testServer.Close()

	address, _ := strings.CutPrefix(testServer.URL, "http")
	wsURL := "ws" + address

	testCases := []csrfTestCase{
		{"http://localhost:8080", "http://localhost:8080", 101},
		{"http://localhost:5000", "http://localhost:8080", http.StatusForbidden},
		{"http://localhost:5000,http://test-domain.local:3000", "http://localhost:8010", http.StatusForbidden},
		{"", "http://localhost:8010", http.StatusForbidden},
		{"", "", http.StatusForbidden},
		{"*", "", 101},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase %v", testCase), func(t *testing.T) {
			os.Setenv("ALLOWED_ORIGINS", testCase.allowedOrigins)

			headers := http.Header{}
			headers.Set("Origin", testCase.origin)

			ws, response, _ := websocket.DefaultDialer.Dial(wsURL, headers)
			if response.StatusCode != testCase.expectedStatusCode {
				t.Fatalf("expected respone status code %d, got %d", testCase.expectedStatusCode, response.StatusCode)
			}

			if ws != nil {
				defer ws.Close()
			}
		})
	}
}
