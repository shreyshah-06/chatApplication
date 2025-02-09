package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"gochatapp/pkg/db"
	"gochatapp/pkg/redisrepo"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	if db.IsUserExist(db.DB, u.Username) {
		jsonResponse(w, false, "Username already taken", nil, 0)
		return
	}

	err := db.RegisterNewUser(db.DB, u.Username, u.Password)
	if err != nil {
		jsonResponse(w, false, "Registration failed", nil, 0)
		return
	}

	jsonResponse(w, true, "User registered successfully", nil, 0)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	err := redisrepo.IsUserAuthentic(u.Username, u.Password)
	if err != nil {
		jsonResponse(w, false, err.Error(), nil, 0)
		return
	}

	jsonResponse(w, true, "Login successful", nil, 0)
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	if !redisrepo.IsUserExist(u.Username) {
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

	if !redisrepo.IsUserExist(u1) || !redisrepo.IsUserExist(u2) {
		jsonResponse(w, false, "Invalid username(s)", nil, 0)
		return
	}

	chats, err := redisrepo.FetchChatBetween(u1, u2, fromTS, toTS)
	if err != nil {
		log.Println("Error fetching chat history:", err)
		jsonResponse(w, false, "Unable to fetch chat history", nil, 0)
		return
	}

	jsonResponse(w, true, "Chat history fetched successfully", chats, len(chats))
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if !redisrepo.IsUserExist(username) {
		jsonResponse(w, false, "Invalid username", nil, 0)
		return
	}

	contacts, err := redisrepo.FetchContactList(username)
	if err != nil {
		log.Println("Error fetching contact list:", err)
		jsonResponse(w, false, "Unable to fetch contact list", nil, 0)
		return
	}

	jsonResponse(w, true, "Contact list fetched successfully", contacts, len(contacts))
}
