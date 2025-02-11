package httpserver

import (
	"encoding/json"
	"net/http"
	"gochatapp/pkg/db"
)


func sendFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	contactUsername := r.URL.Query().Get("contact_username")
	if contactUsername == "" {
		jsonResponse(w, false, "Contact username is required", nil, 0)
		return
	}

	// Validate users
	if !db.IsUserExist(db.DB, u.Username) || !db.IsUserExist(db.DB, contactUsername) {
		jsonResponse(w, false, "Invalid username(s)", nil, 0)
		return
	}

	err := db.SendFollowRequest(db.DB, u.Username, contactUsername)
	if err != nil {
		jsonResponse(w, false, "Failed to send follow request", nil, 0)
		return
	}

	jsonResponse(w, true, "Follow request sent", nil, 0)
}

func acceptFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	contactUsername := r.URL.Query().Get("contact_username")
	if contactUsername == "" {
		jsonResponse(w, false, "Contact username is required", nil, 0)
		return
	}

	// Validate users
	if !db.IsUserExist(db.DB, u.Username) || !db.IsUserExist(db.DB, contactUsername) {
		jsonResponse(w, false, "Invalid username(s)", nil, 0)
		return
	}

	err := db.AcceptFollowRequest(db.DB, u.Username, contactUsername)
	if err != nil {
		jsonResponse(w, false, "Failed to accept follow request", nil, 0)
		return
	}

	jsonResponse(w, true, "Follow request accepted", nil, 0)
}

func rejectFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		jsonResponse(w, false, "Invalid request payload", nil, 0)
		return
	}

	contactUsername := r.URL.Query().Get("contact_username")
	if contactUsername == "" {
		jsonResponse(w, false, "Contact username is required", nil, 0)
		return
	}

	// Validate users
	if !db.IsUserExist(db.DB, u.Username) || !db.IsUserExist(db.DB, contactUsername) {
		jsonResponse(w, false, "Invalid username(s)", nil, 0)
		return
	}

	err := db.RejectFollowRequest(db.DB, u.Username, contactUsername)
	if err != nil {
		jsonResponse(w, false, "Failed to reject follow request", nil, 0)
		return
	}

	jsonResponse(w, true, "Follow request rejected", nil, 0)
}

func pendingFollowRequestsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("u")

	if !db.IsUserExist(db.DB, username) {
		jsonResponse(w, false, "Invalid username", nil, 0)
		return
	}

	requests, err := db.FetchPendingRequests(db.DB, username)
	if err != nil {
		jsonResponse(w, false, "Failed to fetch pending requests", nil, 0)
		return
	}

	jsonResponse(w, true, "Pending requests fetched successfully", requests, len(requests))
}
