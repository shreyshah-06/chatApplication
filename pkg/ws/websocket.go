package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gochatapp/model"
	"gochatapp/pkg/db"
	"gochatapp/pkg/redisrepo"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type string     `json:"type"`
	User string     `json:"user,omitempty"`
	Chat model.Chat `json:"chat,omitempty"`
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *model.Chat)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}
	// register client
	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	receiver(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

// Save chat in PostgreSQL
func saveChatToPostgres(chat *model.Chat) error {
	query := `INSERT INTO chats (from_user, to_user, message, timestamp) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.DB.QueryRow(query, chat.From, chat.To, chat.Msg, chat.Timestamp).Scan(&chat.ID)
	if err != nil {
		return fmt.Errorf("error inserting chat: %v", err)
	}
	return nil
}

// new messages being sent to our WebSocket
// endpoint
func receiver(client *Client) {
	for {
		// Read incoming WebSocket message
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		m := &Message{}
		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("Error unmarshaling chat:", err)
			continue
		}

		fmt.Println("host:", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			// Map client on bootup
			client.Username = m.User
			fmt.Println("Client successfully mapped:", client.Username)
		} else {
			fmt.Println("Received message:", m.Type, m.Chat)
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			// Save in Redis
			id, err := redisrepo.CreateChat(&c)
			if err != nil {
				log.Println("Error saving chat in Redis:", err)
				return
			}
			c.ID = id

			// Save in PostgreSQL
			err = saveChatToPostgres(&c)
			if err != nil {
				log.Println("Error saving chat to PostgreSQL:", err)
				return
			}

			// Broadcast the message
			broadcast <- &c
		}
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	// map our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs)
}

func StartWebsocketServer() {
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	go broadcaster()
	setupRoutes()
	http.ListenAndServe(":8081", nil)
}
