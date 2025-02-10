package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gochatapp/model"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Conn: ws}
	clients[client] = true
	defer func() {
		delete(clients, client)
		client.Conn.Close()
	}()

	fmt.Println("Client connected:", ws.RemoteAddr())
	receiver(client)
}

func receiver(client *Client) {
	for {
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

		if m.Type == "bootup" {
			client.Username = m.User
			fmt.Println("Client successfully mapped:", client.Username)
		} else {
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			// Save in Redis and PostgreSQL via redisrepo.CreateChat
			id, err := redisrepo.CreateChat(&c)
			if err != nil {
				log.Println("Error saving chat:", err)
				return
			}
			c.ID = id

			// Broadcast the message
			broadcast <- &c
		}
	}
}

func Broadcaster() {
	for {
		message := <-broadcast
		for client := range clients {
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