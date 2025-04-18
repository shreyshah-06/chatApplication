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
	Conn       *websocket.Conn
	Username   string
	mu         sync.Mutex
	Registered bool
	// Track previous usernames for switching capability
	PreviousUsernames []string
}

type Message struct {
	Type       string      `json:"type"`
	User       string      `json:"user,omitempty"`
	Chat       *model.Chat `json:"chat,omitempty"`
	Error      string      `json:"error,omitempty"`
	SwitchTo   string      `json:"switch_to,omitempty"`   // New field for identity switching
	SwitchFrom string      `json:"switch_from,omitempty"` // Track previous identity
}

var (
	// Main client map - maps connection pointers to client objects
	clients   = make(map[*Client]bool)
	clientsMu sync.RWMutex

	// Username lookup map - maps usernames to client pointers for quick lookups
	usernameMap   = make(map[string]*Client)
	usernameMu    sync.RWMutex
	
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

// ServeWs handles the initial WebSocket connection
func ServeWs(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	log.Printf("WebSocket connection request from username %s", username)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		Conn:       ws,
		Username:   username,
		Registered: username != "",
		PreviousUsernames: []string{},
	}
	
	log.Printf("New client connected from %s", ws.RemoteAddr())
	if username != "" {
		log.Printf("Username from query param: %s", username)
	}
	
	// Register client in global map
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	// Register in username map if username provided
	if username != "" {
		usernameMu.Lock()
		// If there's an existing client with this username, mark it inactive
		if existingClient, found := usernameMap[username]; found {
			log.Printf("Warning: Username %s is already in use, replacing connection", username)
			// Notify the existing client they're being disconnected
			existingClient.writeJSON(Message{
				Type:  "error",
				Error: "Your session has been taken over by a new connection",
			})
			// Close the existing connection
			existingClient.Conn.Close()
		}
		usernameMap[username] = client
		usernameMu.Unlock()
	}

	// Set initial read deadline
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	
	// Begin handling messages
	go handleClient(client)
}

// handleClient processes all messages for a single client connection
func handleClient(client *Client) {
	// Ensure client cleanup on exit
	defer func() {
		// Clean up client from clients map
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		
		// Clean up client from username map
		if client.Username != "" {
			usernameMu.Lock()
			// Only delete if this client still owns this username
			if currentClient, found := usernameMap[client.Username]; found && currentClient == client {
				delete(usernameMap, client.Username)
			}
			usernameMu.Unlock()
		}
		
		client.Conn.Close()
		log.Printf("Client disconnected: %s", client.Username)
	}()

	// Setup pong handler to keep connection alive
	client.Conn.SetPongHandler(func(string) error {
		return client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	// If we already have a username from the query param, send an immediate acknowledgment
	if client.Username != "" && client.Registered {
		if err := client.writeJSON(Message{
			Type: "ack",
			User: client.Username,
		}); err != nil {
			log.Printf("Error sending initial ack to %s: %v", client.Username, err)
			return
		}
		log.Printf("Auto-registered client with username: %s", client.Username)
	}

	// Main message processing loop
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

		// Parse the incoming message
		var m Message
		if err := json.Unmarshal(p, &m); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			log.Printf("Raw message: %s", string(p))
			client.writeJSON(Message{
				Type: "error",
				Error: "Invalid message format",
			})
			continue
		}

		msgJSON, _ := json.MarshalIndent(m, "", "  ")
		log.Printf("Received message from %s:\n%s", client.Username, msgJSON)
		
		// Handle message based on type
		switch m.Type {
		case "bootup":
			// Handle registration/bootup message
			if m.User == "" {
				log.Printf("Received empty username in bootup")
				client.writeJSON(Message{
					Type: "error",
					Error: "Username cannot be empty",
				})
				continue
			}
			
			// Store previous username if switching
			if client.Username != "" && client.Username != m.User {
				oldUsername := client.Username
				// Add to previous usernames list if not already there
				found := false
				for _, name := range client.PreviousUsernames {
					if name == oldUsername {
						found = true
						break
					}
				}
				if !found {
					client.PreviousUsernames = append(client.PreviousUsernames, oldUsername)
				}
				
				// Remove from username map
				usernameMu.Lock()
				if currentClient, found := usernameMap[oldUsername]; found && currentClient == client {
					delete(usernameMap, oldUsername)
				}
				usernameMu.Unlock()
				
				log.Printf("Client switching from %s to %s", oldUsername, m.User)
			}
			
			// Set new username and registration status
			client.Username = m.User
			client.Registered = true
			
			// Add to username map
			usernameMu.Lock()
			// Handle case where username is already in use
			if existingClient, found := usernameMap[m.User]; found && existingClient != client {
				log.Printf("Warning: Username %s is already in use, replacing connection", m.User)
				// Notify the existing client they're being disconnected
				existingClient.writeJSON(Message{
					Type:  "error",
					Error: "Your session has been taken over by a new connection",
				})
				// Close the existing connection
				existingClient.Conn.Close()
			}
			usernameMap[m.User] = client
			usernameMu.Unlock()
			
			log.Printf("Client registered with username: %s", client.Username)
			
			// Send acknowledgment
			if err := client.writeJSON(Message{
				Type: "ack",
				User: client.Username,
			}); err != nil {
				log.Printf("Error sending ack to %s: %v", client.Username, err)
				return
			}

		case "switch_user":
			// New message type to handle switching between identities
			if m.SwitchTo == "" {
				client.writeJSON(Message{
					Type:  "error",
					Error: "Switch_to username cannot be empty",
				})
				continue
			}
			
			// Store current username in previous list if switching
			if client.Username != "" && client.Username != m.SwitchTo {
				oldUsername := client.Username
				// Add to previous usernames list if not already there
				found := false
				for _, name := range client.PreviousUsernames {
					if name == oldUsername {
						found = true
						break
					}
				}
				if !found {
					client.PreviousUsernames = append(client.PreviousUsernames, oldUsername)
				}
				
				// Remove from username map
				usernameMu.Lock()
				if currentClient, found := usernameMap[oldUsername]; found && currentClient == client {
					delete(usernameMap, oldUsername)
				}
				usernameMu.Unlock()
				
				log.Printf("Client switching from %s to %s", oldUsername, m.SwitchTo)
			}
			
			// Set new username
			client.Username = m.SwitchTo
			client.Registered = true
			
			// Add to username map
			usernameMu.Lock()
			// Handle case where username is already in use
			if existingClient, found := usernameMap[m.SwitchTo]; found && existingClient != client {
				log.Printf("Warning: Username %s is already in use, replacing connection", m.SwitchTo)
				// Notify the existing client they're being disconnected
				existingClient.writeJSON(Message{
					Type:  "error",
					Error: "Your session has been taken over by a new connection",
				})
				// Close the existing connection
				existingClient.Conn.Close()
			}
			usernameMap[m.SwitchTo] = client
			usernameMu.Unlock()
			
			// Send acknowledgment with previous identity info
			if err := client.writeJSON(Message{
				Type:       "switch_ack",
				User:       client.Username,
				SwitchFrom: m.SwitchFrom,
			}); err != nil {
				log.Printf("Error sending switch ack to %s: %v", client.Username, err)
				return
			}
			
			log.Printf("Client switched identity to: %s", client.Username)

		case "chat":
			// Handle chat message
			if !client.Registered || client.Username == "" {
				log.Printf("Received chat message from unregistered client")
				client.writeJSON(Message{
					Type: "error",
					Error: "Please register with bootup message first",
				})
				continue
			}
			
			if m.Chat == nil {
				log.Printf("Received empty chat message from %s", client.Username)
				client.writeJSON(Message{
					Type: "error",
					Error: "Empty message",
				})
				continue
			}

			// Validate message fields
			if m.Chat.From == "" || m.Chat.To == "" || m.Chat.Msg == "" {
				log.Printf("Invalid chat message received from %s", client.Username)
				client.writeJSON(Message{
					Type: "error",
					Error: "Message must include sender, recipient, and content",
				})
				continue
			}

			// FIXED: Respect the specified 'From' field in the message
			// We don't modify the message's From field anymore
			// Just log who's actually sending it for debugging purposes
			if m.Chat.From != client.Username {
				log.Printf("Note: Message from field '%s' doesn't match current username '%s'", 
					m.Chat.From, client.Username)
			}

			// Set timestamp if not already set
			if m.Chat.Timestamp == 0 {
				m.Chat.Timestamp = float64(time.Now().Unix())
			}

			// Save in Redis and get ID
			id, err := redisrepo.CreateChat(m.Chat)
			if err != nil {
				log.Printf("Error saving chat from %s: %v", client.Username, err)
				client.writeJSON(Message{
					Type: "error",
					Error: "Failed to save message",
				})
				continue
			}
			m.Chat.ID = id

			// Broadcast message
			select {
			case broadcast <- m.Chat:
				log.Printf("Message from %s to %s queued for broadcast (ID: %s)", 
					m.Chat.From, m.Chat.To, m.Chat.ID)
				
				// Send immediate confirmation to sender
				client.writeJSON(Message{
					Type: "sent",
					Chat: m.Chat,
				})
				
			default:
				log.Printf("Broadcast channel full, dropped message from %s to %s", 
					m.Chat.From, m.Chat.To)
				client.writeJSON(Message{
					Type: "error",
					Error: "Server busy, please try again",
				})
			}
		}
	}
}

// Broadcaster distributes messages to recipients
func Broadcaster() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-broadcast:
			delivered := false
			
			// Fast lookup for recipient by username
			usernameMu.RLock()
			recipientClient, found := usernameMap[message.To]
			usernameMu.RUnlock()
			
			if found {
				err := recipientClient.writeJSON(Message{
					Type: "chat",
					Chat: message,
				})
				
				if err != nil {
					log.Printf("Error delivering to %s: %v", message.To, err)
					recipientClient.Conn.Close()
					
					// Clean up from maps
					usernameMu.Lock()
					if currentClient, stillExists := usernameMap[message.To]; stillExists && currentClient == recipientClient {
						delete(usernameMap, message.To)
					}
					usernameMu.Unlock()
					
					clientsMu.Lock()
					delete(clients, recipientClient)
					clientsMu.Unlock()
				} else {
					delivered = true
					log.Printf("Successfully delivered message %s to recipient %s", 
						message.ID, message.To)
				}
			}
			
			// If same client is both sender and recipient, make sure they get the message once
			if message.From == message.To {
				delivered = true
				log.Printf("Sender and recipient are the same user: %s", message.From)
			}
			
			if !delivered {
				log.Printf("Recipient %s not connected, message %s stored only", 
					message.To, message.ID)
			}

		case <-ticker.C:
			// Send ping to all clients
			activeCount := 0
			clientsMu.RLock()
			for client := range clients {
				if err := client.Conn.WriteControl(
					websocket.PingMessage,
					[]byte{},
					time.Now().Add(10*time.Second),
				); err != nil {
					log.Printf("Error pinging %s: %v", client.Username, err)
				} else {
					activeCount++
				}
			}
			log.Printf("Ping sent to %d active clients", activeCount)
			clientsMu.RUnlock()
		}
	}
}