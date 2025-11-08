# WebSockets Solution - Deep Dive

## Overview

This solution demonstrates production-grade WebSocket server and client implementation using the gorilla/websocket package. It covers connection management, message broadcasting, room/channel patterns, heartbeat/ping-pong, and graceful shutdown.

## Architecture

### 1. WebSocket Hub Pattern

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    // Client's send buffer is full
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

**Why this pattern:**
- Centralized connection management
- Safe concurrent access to client map
- Non-blocking broadcast (skip slow clients)
- Clean separation of concerns

### 2. Client Connection Handler

```go
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
    id   string
}

func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        // Process message
        c.hub.broadcast <- message
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Hub closed the channel
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**Key design choices:**
- Separate read/write goroutines (one reader, one writer per connection)
- Buffered send channel prevents blocking hub
- Ping/pong for connection health
- Proper deadline handling

### 3. HTTP Upgrade Handler

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // In production, validate origin properly
        return true
    },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
        id:   generateID(),
    }

    client.hub.register <- client

    // Start goroutines
    go client.WritePump()
    go client.ReadPump()
}
```

## Key Patterns

### Pattern 1: Room/Channel System

```go
type Room struct {
    id      string
    clients map[*Client]bool
    mu      sync.RWMutex
}

type RoomHub struct {
    rooms map[string]*Room
    mu    sync.RWMutex
}

func (rh *RoomHub) JoinRoom(roomID string, client *Client) {
    rh.mu.Lock()
    defer rh.mu.Unlock()

    room, ok := rh.rooms[roomID]
    if !ok {
        room = &Room{
            id:      roomID,
            clients: make(map[*Client]bool),
        }
        rh.rooms[roomID] = room
    }

    room.mu.Lock()
    room.clients[client] = true
    room.mu.Unlock()
}

func (rh *RoomHub) LeaveRoom(roomID string, client *Client) {
    rh.mu.RLock()
    room, ok := rh.rooms[roomID]
    rh.mu.RUnlock()

    if !ok {
        return
    }

    room.mu.Lock()
    delete(room.clients, client)
    empty := len(room.clients) == 0
    room.mu.Unlock()

    if empty {
        rh.mu.Lock()
        delete(rh.rooms, roomID)
        rh.mu.Unlock()
    }
}

func (rh *RoomHub) BroadcastToRoom(roomID string, message []byte) {
    rh.mu.RLock()
    room, ok := rh.rooms[roomID]
    rh.mu.RUnlock()

    if !ok {
        return
    }

    room.mu.RLock()
    defer room.mu.RUnlock()

    for client := range room.clients {
        select {
        case client.send <- message:
        default:
            // Skip slow client
        }
    }
}
```

**Use cases:**
- Chat rooms
- Game lobbies
- Collaborative editing
- Live dashboards

### Pattern 2: Message Types and Routing

```go
type Message struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

type MessageHandler interface {
    Handle(client *Client, payload json.RawMessage) error
}

type MessageRouter struct {
    handlers map[string]MessageHandler
    mu       sync.RWMutex
}

func (mr *MessageRouter) Register(msgType string, handler MessageHandler) {
    mr.mu.Lock()
    defer mr.mu.Unlock()
    mr.handlers[msgType] = handler
}

func (mr *MessageRouter) Route(client *Client, data []byte) error {
    var msg Message
    if err := json.Unmarshal(data, &msg); err != nil {
        return err
    }

    mr.mu.RLock()
    handler, ok := mr.handlers[msg.Type]
    mr.mu.RUnlock()

    if !ok {
        return fmt.Errorf("unknown message type: %s", msg.Type)
    }

    return handler.Handle(client, msg.Payload)
}

// Example handlers
type ChatMessageHandler struct {
    hub *Hub
}

func (h *ChatMessageHandler) Handle(client *Client, payload json.RawMessage) error {
    var chatMsg struct {
        Text string `json:"text"`
    }

    if err := json.Unmarshal(payload, &chatMsg); err != nil {
        return err
    }

    // Broadcast to all clients
    response, _ := json.Marshal(Message{
        Type: "chat",
        Payload: json.RawMessage(fmt.Sprintf(`{"user":"%s","text":"%s"}`,
            client.id, chatMsg.Text)),
    })

    h.hub.broadcast <- response
    return nil
}
```

### Pattern 3: Connection Health Monitoring

```go
const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512 * 1024 // 512 KB
)

func (c *Client) SetupHealthCheck() {
    // Set max message size
    c.conn.SetReadLimit(maxMessageSize)

    // Set read deadline
    c.conn.SetReadDeadline(time.Now().Add(pongWait))

    // Set pong handler
    c.conn.SetPongHandler(func(appData string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return // Connection dead
            }
        // ... handle messages
        }
    }
}
```

**Why ping/pong:**
- Detect dead connections
- Keep connections alive through proxies/firewalls
- Prevent resource leaks

### Pattern 4: Graceful Shutdown

```go
type Server struct {
    hub        *Hub
    httpServer *http.Server
    shutdown   chan struct{}
}

func (s *Server) Shutdown(ctx context.Context) error {
    close(s.shutdown)

    // Stop accepting new connections
    if err := s.httpServer.Shutdown(ctx); err != nil {
        return err
    }

    // Close all WebSocket connections gracefully
    s.hub.mu.Lock()
    for client := range s.hub.clients {
        client.conn.WriteMessage(
            websocket.CloseMessage,
            websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down"),
        )
        client.conn.Close()
    }
    s.hub.mu.Unlock()

    return nil
}

// Usage in main.go
func main() {
    server := NewServer()

    go server.Start()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    // Shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Performance Optimization

### 1. Connection Pooling

```go
type ClientPool struct {
    pool sync.Pool
}

func NewClientPool() *ClientPool {
    return &ClientPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &Client{
                    send: make(chan []byte, 256),
                }
            },
        },
    }
}

func (cp *ClientPool) Get() *Client {
    return cp.pool.Get().(*Client)
}

func (cp *ClientPool) Put(client *Client) {
    // Reset client state
    client.conn = nil
    client.hub = nil
    client.id = ""

    // Return to pool
    cp.pool.Put(client)
}
```

### 2. Message Batching

```go
type BatchWriter struct {
    client   *Client
    buffer   [][]byte
    mu       sync.Mutex
    ticker   *time.Ticker
    maxBatch int
}

func (bw *BatchWriter) Write(message []byte) {
    bw.mu.Lock()
    defer bw.mu.Unlock()

    bw.buffer = append(bw.buffer, message)

    if len(bw.buffer) >= bw.maxBatch {
        bw.flush()
    }
}

func (bw *BatchWriter) flush() {
    if len(bw.buffer) == 0 {
        return
    }

    // Combine messages
    combined := bytes.Join(bw.buffer, []byte("\n"))

    bw.client.send <- combined
    bw.buffer = bw.buffer[:0]
}

func (bw *BatchWriter) Start() {
    bw.ticker = time.NewTicker(100 * time.Millisecond)

    go func() {
        for range bw.ticker.C {
            bw.mu.Lock()
            bw.flush()
            bw.mu.Unlock()
        }
    }()
}
```

### 3. Binary Protocol

```go
// Instead of JSON, use binary encoding for performance
type BinaryMessage struct {
    Type    uint8
    Length  uint32
    Payload []byte
}

func (bm *BinaryMessage) Encode() []byte {
    buf := make([]byte, 5+len(bm.Payload))
    buf[0] = bm.Type
    binary.BigEndian.PutUint32(buf[1:5], bm.Length)
    copy(buf[5:], bm.Payload)
    return buf
}

func DecodeBinaryMessage(data []byte) (*BinaryMessage, error) {
    if len(data) < 5 {
        return nil, errors.New("message too short")
    }

    return &BinaryMessage{
        Type:    data[0],
        Length:  binary.BigEndian.Uint32(data[1:5]),
        Payload: data[5:],
    }, nil
}
```

## Security Considerations

### 1. Origin Validation

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")

        // Allow same origin
        if origin == "" {
            return true
        }

        // Validate against whitelist
        allowedOrigins := []string{
            "https://example.com",
            "https://app.example.com",
        }

        for _, allowed := range allowedOrigins {
            if origin == allowed {
                return true
            }
        }

        return false
    },
}
```

### 2. Authentication

```go
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // Extract token from query param or header
    token := r.URL.Query().Get("token")
    if token == "" {
        token = r.Header.Get("Authorization")
    }

    // Validate token
    userID, err := validateToken(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    client := &Client{
        hub:    hub,
        conn:   conn,
        send:   make(chan []byte, 256),
        id:     generateID(),
        userID: userID, // Associate with authenticated user
    }

    // ... register client
}
```

### 3. Rate Limiting

```go
type RateLimiter struct {
    clients map[string]*rate.Limiter
    mu      sync.RWMutex
    rate    rate.Limit
    burst   int
}

func (rl *RateLimiter) GetLimiter(clientID string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.clients[clientID]
    rl.mu.RUnlock()

    if !exists {
        rl.mu.Lock()
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.clients[clientID] = limiter
        rl.mu.Unlock()
    }

    return limiter
}

func (c *Client) ReadPump() {
    rateLimiter := NewRateLimiter(rate.Limit(10), 20) // 10 msg/sec, burst 20

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        // Check rate limit
        limiter := rateLimiter.GetLimiter(c.id)
        if !limiter.Allow() {
            c.send <- []byte(`{"error":"rate limit exceeded"}`)
            continue
        }

        // Process message
        c.hub.broadcast <- message
    }
}
```

### 4. Input Validation

```go
const maxMessageSize = 512 * 1024 // 512 KB

func (c *Client) ReadPump() {
    c.conn.SetReadLimit(maxMessageSize)

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        // Validate message structure
        if !json.Valid(message) {
            c.send <- []byte(`{"error":"invalid JSON"}`)
            continue
        }

        // Validate message content
        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            c.send <- []byte(`{"error":"malformed message"}`)
            continue
        }

        // Process validated message
        c.router.Route(c, message)
    }
}
```

## Testing Strategies

### 1. WebSocket Client Test

```go
func TestWebSocket(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ServeWS(hub, w, r)
    }))
    defer server.Close()

    // Convert http:// to ws://
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

    // Connect client
    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    require.NoError(t, err)
    defer ws.Close()

    // Send message
    err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
    require.NoError(t, err)

    // Read response
    _, msg, err := ws.ReadMessage()
    require.NoError(t, err)
    assert.Equal(t, "hello", string(msg))
}
```

### 2. Hub Test

```go
func TestHub(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    // Create mock clients
    client1 := &Client{send: make(chan []byte, 10)}
    client2 := &Client{send: make(chan []byte, 10)}

    // Register clients
    hub.register <- client1
    hub.register <- client2

    time.Sleep(10 * time.Millisecond) // Let registration process

    // Broadcast message
    hub.broadcast <- []byte("test message")

    time.Sleep(10 * time.Millisecond) // Let broadcast process

    // Verify both clients received message
    select {
    case msg := <-client1.send:
        assert.Equal(t, "test message", string(msg))
    case <-time.After(time.Second):
        t.Fatal("client1 did not receive message")
    }

    select {
    case msg := <-client2.send:
        assert.Equal(t, "test message", string(msg))
    case <-time.After(time.Second):
        t.Fatal("client2 did not receive message")
    }
}
```

## Common Pitfalls

### 1. Concurrent Write

```go
// BAD: Multiple goroutines writing to same connection
go func() {
    conn.WriteMessage(websocket.TextMessage, msg1)
}()
go func() {
    conn.WriteMessage(websocket.TextMessage, msg2) // RACE!
}()

// GOOD: Single writer goroutine
func (c *Client) WritePump() {
    for message := range c.send {
        c.conn.WriteMessage(websocket.TextMessage, message)
    }
}
```

### 2. Not Closing Connections

```go
// BAD: Connection leak
func ServeWS(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    // conn never closed!
}

// GOOD: Ensure cleanup
func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close() // Always close
    }()
    // ...
}
```

### 3. Blocking Broadcast

```go
// BAD: Blocks if client is slow
for client := range h.clients {
    client.send <- message // Can block!
}

// GOOD: Non-blocking send
for client := range h.clients {
    select {
    case client.send <- message:
    default:
        // Skip slow client or close it
    }
}
```

## Production Checklist

- [ ] Ping/pong health checks implemented
- [ ] Graceful shutdown handling
- [ ] Origin validation in production
- [ ] Authentication/authorization
- [ ] Rate limiting per client
- [ ] Message size limits
- [ ] Connection limits per IP/user
- [ ] Proper error handling and logging
- [ ] Metrics (active connections, messages/sec)
- [ ] Separate read/write goroutines
- [ ] Channel buffer sizes tuned
- [ ] Memory profiling done
- [ ] Load testing completed

## Further Reading

- **gorilla/websocket:** https://github.com/gorilla/websocket
- **WebSocket RFC:** https://datatracker.ietf.org/doc/html/rfc6455
- **Chat example:** https://github.com/gorilla/websocket/tree/master/examples/chat
- **nhooyr.io/websocket:** https://github.com/nhooyr/websocket (alternative library)
