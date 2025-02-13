package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"gochatapp/pkg/db"
	"gochatapp/pkg/redisrepo"
	"gochatapp/utils"
)

type userInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userReq struct {
	Username string `json:"username"`
}

type response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func jsonResponse(w http.ResponseWriter, status bool, message string, data interface{}, total int) {
	setJSONHeader(w)
	json.NewEncoder(w).Encode(response{
		Status:  status,
		Message: message,
		Data:    data,
		Total:   total,
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	u := &userInfo{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	if db.IsUserExist(db.DB, u.Username) {
		jsonResponse(w, false, "Username already taken", nil, 0)
		return
	}

	// Hash the user's password before storing it
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		// If there is an error hashing the password, return an error response
		jsonResponse(w, false, "Error hashing password", nil, 0)
		return
	}

	// Store the user with the hashed password in the database
	err = db.RegisterNewUser(db.DB, u.Username, hashedPassword)
	if err != nil {
		// If the registration fails, return an error response
		jsonResponse(w, false, "Registration failed", nil, 0)
		return
	}

	// If everything is successful, return a success response
	jsonResponse(w, true, "User registered successfully", nil, 0)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	u := &userInfo{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	err := db.IsUserAuthentic(db.DB, u.Username, u.Password)
	if err != nil {
		jsonResponse(w, false, err.Error(), nil, 0)
		return
	}
	// If authentication is successful, generate the JWT token
	token, err := utils.CreateJWT(u.Username)
	if err != nil {
		// If there is an error generating the token, return an error response
		jsonResponse(w, false, "Error generating JWT token", nil, 0)
		return
	}

	// Return a successful response with the JWT token
	jsonResponse(w, true, "Login successful", map[string]string{"token": token}, 0)
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	if !db.IsUserExist(db.DB ,u.Username) {
		jsonResponse(w, false, "Invalid username", nil, 0)
		return
	}

	jsonResponse(w, true, "Contact verified", nil, 0)
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")
	fromTS := r.URL.Query().Get("from-ts")
	toTS := r.URL.Query().Get("to-ts")

	if fromTS == "" {
		fromTS = "0"
	}
	if toTS == "" {
		toTS = "+inf"
	}

	if !db.IsUserExist(db.DB ,u1) || !db.IsUserExist(db.DB ,u2) {
		jsonResponse(w, false, "Invalid username(s)", nil, 0)
		return
	}

	// Try fetching recent chat from Redis
	chats, err := redisrepo.FetchChatBetween(u1, u2, fromTS, toTS)
	if err != nil {
		log.Println("Error fetching chat history from Redis:", err)
	}

	// If no chat found in Redis, try fetching from PostgreSQL
	if len(chats) == 0 {
		chats, err = db.FetchChatBetween(u1, u2, fromTS, toTS)
		if err != nil {
			log.Println("Error fetching chat history from PostgreSQL:", err)
			jsonResponse(w, false, "Unable to fetch chat history", nil, 0)
			return
		}
	}

	// Return the chat history
	jsonResponse(w, true, "Chat history fetched successfully", chats, len(chats))
}


func contactListHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if !db.IsUserExist(db.DB, username) {
		jsonResponse(w, false, "Invalid username", nil, 0)
		return
	}

	contacts, err := db.FetchContactList(db.DB, username)
	if err != nil {
		log.Println("Error fetching contact list:", err)
		jsonResponse(w, false, "Unable to fetch contact list", nil, 0)
		return
	}

	jsonResponse(w, true, "Contact list fetched successfully", contacts, len(contacts))
}
