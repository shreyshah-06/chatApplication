package httpserver

import (
	"fmt"
	"gochatapp/pkg/redisrepo"
	"gochatapp/pkg/ws"
	"net/http"

	// "gochatapp/pkg/redisrepo"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// StartHTTPServer initializes the HTTP server
func StartHTTPServer() {
	// Initialize Redis connection
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	// Create necessary indexes
	redisrepo.CreateFetchChatBetweenIndex()

	// Define routes
	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler).Methods(http.MethodGet)
	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodGet)

	r.HandleFunc("/ws", ws.ServeWs)

	// Start server with CORS configuration
	handler := cors.AllowAll().Handler(r)
	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", handler)
}

// statusHandler handles server status checks
func statusHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, true, "Server is running", nil, 0)
}
