package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	parallelConnections := flag.Int("n", 1, "Number of parallel connections to open")
	serverAddr := flag.String("server", "localhost:8080", "server address")
	flag.Parse()

	if *parallelConnections < 1 {
		log.Println("No connections to open. Exiting...")
		return
	}

	wsURL := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/goapp/ws"}
	fmt.Printf("Connecting to ws server: %s\n", wsURL.String())

	var wg sync.WaitGroup
	wg.Add(*parallelConnections)

	for i := 0; i < *parallelConnections; i++ {
		go func(connID int) {
			defer wg.Done()
			connectWebSocket(connID, wsURL.String())
		}(i)
	}
	wg.Wait()
}

func connectWebSocket(connID int, serverURL string) {
	ws, response, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Printf("[conn #%d] Failed to connect: %v\n", connID, err)
		return
	}
	if response.StatusCode != 101 {
		log.Printf("[conn #%d] Can't upgrade http connection. Status code %d\n", connID, response.StatusCode)
		return
	}
	defer ws.Close()

	log.Printf("[conn #%d] Connected to server\n", connID)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[conn #%d] Connection closed: %v\n", connID, err)
			return
		}
		log.Printf("[conn #%d] %s\n", connID, string(message))
	}
}
