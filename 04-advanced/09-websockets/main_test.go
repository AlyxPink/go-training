package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub() returned nil")
	}
	if hub.clients == nil {
		t.Error("hub.clients should be initialized")
	}
	if hub.broadcast == nil {
		t.Error("hub.broadcast should be initialized")
	}
	if hub.register == nil {
		t.Error("hub.register should be initialized")
	}
	if hub.unregister == nil {
		t.Error("hub.unregister should be initialized")
	}
}

func TestHubRegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock client
	client := &Client{
		send: make(chan []byte, 256),
		hub:  hub,
		id:   "test-client",
	}

	// Register the client
	hub.register <- client

	// Give the hub time to process
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if !hub.clients[client] {
		t.Error("Client was not registered with hub")
	}
}

func TestHubUnregisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create and register a client
	client := &Client{
		send: make(chan []byte, 256),
		hub:  hub,
		id:   "test-client",
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Unregister the client
	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if hub.clients[client] {
		t.Error("Client was not unregistered from hub")
	}
}

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create and register two clients
	client1 := &Client{
		send: make(chan []byte, 256),
		hub:  hub,
		id:   "client-1",
	}
	client2 := &Client{
		send: make(chan []byte, 256),
		hub:  hub,
		id:   "client-2",
	}

	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	// Broadcast a message
	testMessage := []byte("test broadcast")
	hub.broadcast <- testMessage
	time.Sleep(10 * time.Millisecond)

	// Check if both clients received the message
	select {
	case msg := <-client1.send:
		if string(msg) != string(testMessage) {
			t.Errorf("Client1 received wrong message: got %s, want %s", msg, testMessage)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive broadcast message")
	}

	select {
	case msg := <-client2.send:
		if string(msg) != string(testMessage) {
			t.Errorf("Client2 received wrong message: got %s, want %s", msg, testMessage)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 did not receive broadcast message")
	}
}

func TestServeWs(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect as a WebSocket client
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket: %v", err)
	}
	defer conn.Close()

	// Give time for registration
	time.Sleep(50 * time.Millisecond)

	// Verify client was registered
	hub.mu.RLock()
	clientCount := len(hub.clients)
	hub.mu.RUnlock()

	if clientCount != 1 {
		t.Errorf("Expected 1 client, got %d", clientCount)
	}

	// Test sending a message
	testMsg := "Hello WebSocket"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(testMsg)); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Read the echoed message back
	conn.SetReadDeadline(time.Now().Add(time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	if string(msg) != testMsg {
		t.Errorf("Expected message %q, got %q", testMsg, msg)
	}
}

func TestServeHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	serveHome(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "WebSocket") {
		t.Error("Home page should contain 'WebSocket'")
	}
}

func TestServeHomeNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	w := httptest.NewRecorder()

	serveHome(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestServeHomeMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	serveHome(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
