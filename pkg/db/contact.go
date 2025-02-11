package db

import (
	"database/sql"
	"gochatapp/model"
	"log"
)

// SendFollowRequest adds a follow request in the contacts table
func SendFollowRequest(db *sql.DB, username, contactUsername string) error {
	query := `
		INSERT INTO contacts (username, contact_username, status) 
		VALUES ($1, $2, 'pending') 
		ON CONFLICT (username, contact_username) DO NOTHING;
	`
	_, err := db.Exec(query, username, contactUsername)
	if err != nil {
		log.Println("Error sending follow request:", err)
	}
	return err
}

// AcceptFollowRequest updates the status of a follow request to 'accepted'
func AcceptFollowRequest(db *sql.DB, username, contactUsername string) error {
	query := `
		UPDATE contacts 
		SET status = 'accepted', updated_at = NOW() 
		WHERE username = $1 AND contact_username = $2 AND status = 'pending';
	`
	_, err := db.Exec(query, contactUsername, username)
	if err != nil {
		log.Println("Error accepting follow request:", err)
	}
	return err
}

// RejectFollowRequest updates the status of a follow request to 'rejected'
func RejectFollowRequest(db *sql.DB, username, contactUsername string) error {
	query := `
		UPDATE contacts 
		SET status = 'rejected', updated_at = NOW() 
		WHERE username = $1 AND contact_username = $2 AND status = 'pending';
	`
	_, err := db.Exec(query, contactUsername, username)
	if err != nil {
		log.Println("Error rejecting follow request:", err)
	}
	return err
}

// FetchContactList fetches accepted contacts for a user
func FetchContactList(db *sql.DB, username string) ([]model.ContactList, error) {
	query := `
		SELECT contact_username 
		FROM contacts 
		WHERE username = $1 AND status = 'accepted'
		ORDER BY updated_at DESC;
	`

	rows, err := db.Query(query, username)
	if err != nil {
		log.Println("Error fetching contact list from PostgreSQL:", err)
		return nil, err
	}
	defer rows.Close()

	var contacts []model.ContactList
	for rows.Next() {
		var contact model.ContactList
		if err := rows.Scan(&contact.Username); err != nil {
			log.Println("Error scanning contact list row:", err)
			continue
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// FetchPendingRequests fetches pending follow requests for a user
func FetchPendingRequests(db *sql.DB, username string) ([]model.ContactList, error) {
	query := `
		SELECT contact_username AS username, EXTRACT(EPOCH FROM created_at) AS last_activity
		FROM contacts 
		WHERE username = $1 AND status = 'pending'
		ORDER BY created_at DESC;
	`
	rows, err := db.Query(query, username)
	if err != nil {
		log.Println("Error fetching pending requests from PostgreSQL:", err)
		return nil, err
	}
	defer rows.Close()

	var pendingRequests []model.ContactList
	for rows.Next() {
		var request model.ContactList
		if err := rows.Scan(&request.Username, &request.LastActivity); err != nil {
			log.Println("Error scanning pending request row:", err)
			continue
		}
		pendingRequests = append(pendingRequests, request)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	return pendingRequests, nil
}
