package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"gochatapp/model"
	"gochatapp/pkg/db"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func RegisterNewUser(username, password string) error {
	err := redisClient.Set(context.Background(), username, password, 0).Err()
	if err != nil {
		log.Println("Error while adding new user:", err)
		return err
	}

	err = redisClient.SAdd(context.Background(), userSetKey(), username).Err()
	if err != nil {
		log.Println("Error while adding user in set:", err)
		redisClient.Del(context.Background(), username)
		return err
	}
	return nil
}

func IsUserExist(username string) bool {
	return redisClient.SIsMember(context.Background(), userSetKey(), username).Val()
}

func IsUserAuthentic(username, password string) error {
	p := redisClient.Get(context.Background(), username).Val()
	if !strings.EqualFold(p, password) {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}

func UpdateContactList(username, contact string) error {
	zs := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}
	err := redisClient.ZAdd(context.Background(), contactListZKey(username), zs).Err()
	if err != nil {
		log.Println("Error updating contact list for user:", username, "Contact:", contact, err)
		return err
	}
	return nil
}

func CreateChat(c *model.Chat) (string, error) {
	chatKey := chatKey()
	
	// Serialize the chat object
	by, err := json.Marshal(c)
	if err != nil {
		log.Println("Error marshaling chat:", err)
		return "", err
	}
	
	// Save chat in Redis concurrently using standard SET command instead of JSON.SET
	redisDone := make(chan string, 1)
	go func() {
		res, err := redisClient.Set(
			context.Background(),
			chatKey,
			string(by),
			0, // no expiration
		).Result()
		
		if err != nil {
			log.Println("Error setting chat in Redis:", err)
		} else {
			log.Println("Chat successfully set in Redis:", res)
			redisDone <- chatKey
		}
		close(redisDone)
	}()
	
	// Update contact lists concurrently
	go func() {
		if err := UpdateContactList(c.From, c.To); err != nil {
			log.Println("Error updating contact list for", c.From)
		}
		if err := UpdateContactList(c.To, c.From); err != nil {
			log.Println("Error updating contact list for", c.To)
		}
	}()
	
	// Store chat in PostgreSQL concurrently
	go func() {
		if err := db.StoreChatInPostgres(c); err != nil {
			log.Println("Error storing chat in PostgreSQL:", err)
		}
	}()
	
	// Wait for Redis operation to complete to return key
	chatKey = <-redisDone
	return chatKey, nil
}



func CreateFetchChatBetweenIndex() {
	res, err := redisClient.Do(context.Background(),
		"FT.CREATE", chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA",
		"$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	if err != nil {
		log.Println("Error creating chat index:", err)
		return
	}

	fmt.Println("Index created successfully:", res)
}

func FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	// Fetch messages from Redis with a limit (20 for recent messages)
	limit := "20" // or pass the limit dynamically based on your server load or user preferences

	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		username1, username2, username1, username2, fromTS, toTS)

	// Redis search query to fetch messages in batches (20 messages at a time)
	res, err := redisClient.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
		"LIMIT", "0", limit,
	).Result()

	if err != nil {
		log.Println("Error fetching chat between users:", err)
		return nil, err
	}

	// Deserialize the response from Redis
	data := Deserialise(res)
	if data == nil {
		return nil, fmt.Errorf("no data found for the given query")
	}

	// Convert deserialized data to model.Chat objects
	return DeserialiseChat(data), nil
}

func FetchContactList(username string) ([]model.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(username),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}

	res, err := redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()
	if err != nil {
		log.Println("Error fetching contact list for user:", username, err)
		return nil, err
	}

	return DeserialiseContactList(res), nil
}
