package httpserver

import (
	"fmt"
	auth "gochatapp/pkg/middleware"
	"gochatapp/pkg/redisrepo"
	"gochatapp/pkg/ws"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// StartHTTPServer initializes the HTTP server and defines the routes
func StartHTTPServer() {
	// Initialize Redis connection
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close() // Ensure Redis connection is closed after the server shuts down

	// Create necessary indexes for Redis (e.g., for chat history)
	redisrepo.CreateFetchChatBetweenIndex()

	// Create a new router
	r := mux.NewRouter()

	// Server status route (for health check)
	r.HandleFunc("/status", statusHandler).Methods(http.MethodGet)

	// Authentication routes
	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost) // User registration route
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)       // User login route

	// Protected routes with JWT authentication middleware
	r.Handle("/verify-contact", auth.JwtMiddleware(http.HandlerFunc(verifyContactHandler))).Methods(http.MethodPost)
	r.Handle("/chat-history", auth.JwtMiddleware(http.HandlerFunc(chatHistoryHandler))).Methods(http.MethodGet)

	// Contact management routes (user interaction related to contacts)
	r.Handle("/contact-list", auth.JwtMiddleware(http.HandlerFunc(contactListHandler))).Methods(http.MethodGet)
	r.Handle("/send-follow-request", auth.JwtMiddleware(http.HandlerFunc(sendFollowRequestHandler))).Methods(http.MethodPost)
	r.Handle("/accept-follow-request", auth.JwtMiddleware(http.HandlerFunc(acceptFollowRequestHandler))).Methods(http.MethodPut)
	r.Handle("/reject-follow-request", auth.JwtMiddleware(http.HandlerFunc(rejectFollowRequestHandler))).Methods(http.MethodPut)
	r.Handle("/pending-follow-request", auth.JwtMiddleware(http.HandlerFunc(pendingFollowRequestsHandler))).Methods(http.MethodGet)

	// WebSocket route for real-time communication
	r.Handle("/ws", auth.JwtMiddleware(http.HandlerFunc(ws.ServeWs)))

	// Start the server with CORS configuration (Allow all origins for simplicity, can be restricted as needed)
	handler := cors.AllowAll().Handler(r)

	// Print server start message and begin listening on port 8080
	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", handler)
}

// statusHandler handles server status checks
func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Respond with a JSON message confirming the server is running
	jsonResponse(w, true, "Server is running", nil, 0)
}
