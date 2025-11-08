package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in this example
	},
}

// Client represents a WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	id   string
}

// Hub manages all active clients and broadcasts
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	// TODO: Initialize and return Hub
	return nil
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	// TODO: Implement hub's main event loop
	// - Listen on register channel
	// - Listen on unregister channel
	// - Listen on broadcast channel
	// - Handle each event appropriately
}

// readPump handles incoming messages from client
func (c *Client) readPump() {
	// TODO: Implement read loop
	// - Read messages from WebSocket
	// - Send to hub.broadcast
	// - Handle errors and disconnection
	// - Unregister client on exit
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
}

// writePump handles outgoing messages to client
func (c *Client) writePump() {
	// TODO: Implement write loop
	// - Read from c.send channel
	// - Write to WebSocket connection
	// - Handle errors
	// - Close connection on exit
	defer c.conn.Close()
}

// serveWs handles WebSocket upgrade and client creation
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Upgrade HTTP connection to WebSocket
	// TODO: Create new Client
	// TODO: Register client with hub
	// TODO: Start readPump and writePump goroutines
}

// serveHome serves the HTML page
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	addr := ":8080"
	fmt.Printf("WebSocket server starting on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
