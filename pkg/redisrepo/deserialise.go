package redisrepo

import (
	"encoding/json"
	"gochatapp/model"
	"log"

	"github.com/go-redis/redis/v8"
)

type Document struct {
	ID      string `json:"_id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}

func Deserialise(res interface{}) []Document {
	var docs []Document
	switch v := res.(type) {
	case []interface{}:
		if len(v) > 1 {
			total := len(v) - 1
			docs = make([]Document, 0, total/2)
			for i := 1; i <= total; i += 2 {
				arr := v[i+1].([]interface{})
				value := arr[len(arr)-1].(string)
				docs = append(docs, Document{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				})
			}
		}
	default:
		log.Printf("Unexpected response type: %T", res)
	}
	return docs
}

func DeserialiseChat(docs []Document) []model.Chat {
	chats := make([]model.Chat, len(docs))
	for i, doc := range docs {
		var c model.Chat
		if err := json.Unmarshal(doc.Payload, &c); err == nil {
			c.ID = doc.ID
			chats[i] = c
		}
	}
	return chats
}

func DeserialiseContactList(contacts []redis.Z) []model.ContactList {
	contactList := make([]model.ContactList, 0, len(contacts))
	for _, contact := range contacts {
		contactList = append(contactList, model.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}
	return contactList
}