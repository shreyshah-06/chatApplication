package ws

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "gochatapp/model"
    "gochatapp/pkg/redisrepo"

    "github.com/gorilla/websocket"
)

type Client struct {
    Conn     *websocket.Conn
    Username string
    mu       sync.Mutex
}

type Message struct {
    Type string      `json:"type"`
    User string      `json:"user,omitempty"`
    Chat *model.Chat `json:"chat,omitempty"`
}

var (
    clients   = make(map[*Client]bool)
    clientsMu sync.RWMutex
    broadcast = make(chan *model.Chat, 256)
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Client) writeJSON(v interface{}) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.Conn.WriteJSON(v)
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }

    client := &Client{Conn: ws}
    
    clientsMu.Lock()
    clients[client] = true
    clientsMu.Unlock()

    defer func() {
        clientsMu.Lock()
        delete(clients, client)
        clientsMu.Unlock()
        client.Conn.Close()
        log.Printf("Client disconnected: %s", client.Username)
    }()

    // Set initial read deadline
    ws.SetReadDeadline(time.Now().Add(60 * time.Second))
    
    log.Printf("New client connected: %s", ws.RemoteAddr())
    receiver(client)
}

func receiver(client *Client) {
    client.Conn.SetPongHandler(func(string) error {
        return client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    })

    for {
        _, p, err := client.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Unexpected closure for %s: %v", client.Username, err)
            }
            return
        }

        // Reset read deadline after successful read
        client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

        var m Message
        if err := json.Unmarshal(p, &m); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        switch m.Type {
        case "bootup":
            client.Username = m.User
            log.Printf("Client registered: %s", client.Username)
            
            // Send acknowledgment
            if err := client.writeJSON(Message{
                Type: "ack",
                User: client.Username,
            }); err != nil {
                log.Printf("Error sending ack to %s: %v", client.Username, err)
                return
            }

        case "chat":
            if m.Chat == nil {
                log.Printf("Received empty chat message from %s", client.Username)
                continue
            }

            // Validate message fields
            if m.Chat.From == "" || m.Chat.To == "" || m.Chat.Msg == "" {
                log.Printf("Invalid chat message received from %s", client.Username)
                continue
            }

            // Ensure the sender matches the connected client
            if m.Chat.From != client.Username {
                log.Printf("Message sender mismatch. Expected %s, got %s", client.Username, m.Chat.From)
                continue
            }

            // Set timestamp if not already set
            if m.Chat.Timestamp == 0 {
                m.Chat.Timestamp = float64(time.Now().Unix())
            }

            // Save in Redis and get ID
            id, err := redisrepo.CreateChat(m.Chat)
            if err != nil {
                log.Printf("Error saving chat from %s: %v", client.Username, err)
                continue
            }
            m.Chat.ID = id

            // Broadcast message
            select {
            case broadcast <- m.Chat:
                log.Printf("Message from %s to %s for broadcast", m.Chat.From, m.Chat.To)
            default:
                log.Printf("Broadcast channel full, dropped message from %s to %s", m.Chat.From, m.Chat.To)
            }
        }
    }
}

func Broadcaster() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case message := <-broadcast:
            clientsMu.RLock()
            for client := range clients {
                if client.Username == message.From || client.Username == message.To {
                    err := client.writeJSON(Message{
                        Type: "chat",
                        Chat: message,
                    })
                    
                    if err != nil {
                        log.Printf("Error broadcasting to %s: %v", client.Username, err)
                        
                        clientsMu.RUnlock()
                        clientsMu.Lock()
                        delete(clients, client)
                        clientsMu.Unlock()
                        clientsMu.RLock()
                        
                        client.Conn.Close()
                        continue
                    }
                    
                    log.Printf("Successfully delivered message %s to %s", message.ID, client.Username)
                }
            }
            clientsMu.RUnlock()

        case <-ticker.C:
            clientsMu.RLock()
            for client := range clients {
                if err := client.Conn.WriteControl(
                    websocket.PingMessage,
                    []byte{},
                    time.Now().Add(10*time.Second),
                ); err != nil {
                    log.Printf("Error pinging %s: %v", client.Username, err)
                }
            }
            clientsMu.RUnlock()
        }
    }
}