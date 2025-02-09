package redisrepo

import (
	"fmt"
	"time"
)

func userSetKey() string {
	return "users"
}

// sessionKey generates a key for client sessions.
// Currently unused but can be useful for managing session data in Redis.
// func sessionKey(client string) string {
// 	return "session#" + client
// }

func chatKey() string {
	return fmt.Sprintf("chat#%d", time.Now().UnixMilli())
}

func chatIndex() string {
	return "idx#chats"
}

func contactListZKey(username string) string {
	return "contacts:" + username
}
