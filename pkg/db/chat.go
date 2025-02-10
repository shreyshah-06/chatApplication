package db

import (
	"fmt"
	"gochatapp/model"
	"log"
)

// StoreChatInPostgres stores the chat data in PostgreSQL
func StoreChatInPostgres(c *model.Chat) error {
	fmt.Println(c)
	query := `INSERT INTO messages (sender, receiver, content, sent_at)
          VALUES ($1, $2, $3, to_timestamp($4)) RETURNING id`


	// Execute the query and retrieve the generated ID
	err := DB.QueryRow(query, c.From, c.To, c.Msg, c.Timestamp).Scan(&c.ID)
	if err != nil {
		log.Println("Error storing chat in PostgreSQL:", err)
		return err
	}

	return nil
}

// FetchChatBetween retrieves chat history between two users from PostgreSQL
func FetchChatBetween(u1, u2, fromTS, toTS string) ([]model.Chat, error) {
	var chats []model.Chat

	query := `SELECT id, sender, receiver, content, extract(epoch from sent_at) as timestamp
				FROM messages
				WHERE (sender = $1 AND receiver = $2 OR sender = $2 AND receiver = $1)
				AND sent_at BETWEEN to_timestamp($3) AND to_timestamp($4)
				ORDER BY sent_at DESC`

	rows, err := DB.Query(query, u1, u2, fromTS, toTS)
	if err != nil {
		log.Println("Error fetching chat history from PostgreSQL:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		// Scan the result into the Chat struct
		if err := rows.Scan(&chat.ID, &chat.From, &chat.To, &chat.Msg, &chat.Timestamp); err != nil {
			log.Println("Error scanning chat data:", err)
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func FetchOldChatBetween(u1, u2, fromTS, toTS, limit string) ([]model.Chat, error) {
	var chats []model.Chat

	// PostgreSQL query to fetch older chat history based on timestamp
	query := `SELECT id, sender, receiver, content, extract(epoch from sent_at) as timestamp
				FROM messages
				WHERE (sender = $1 AND receiver = $2 OR sender = $2 AND receiver = $1)
				AND sent_at < to_timestamp($3)
				ORDER BY sent_at DESC
				LIMIT $4`

	rows, err := DB.Query(query, u1, u2, fromTS, limit)
	if err != nil {
		log.Println("Error fetching older chat history from PostgreSQL:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		// Scan the result into the Chat struct
		if err := rows.Scan(&chat.ID, &chat.From, &chat.To, &chat.Msg, &chat.Timestamp); err != nil {
			log.Println("Error scanning chat data:", err)
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
