package model

type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	Timestamp float64  `json:"timestamp"`
}

type ContactList struct {
	Username     string `json:"username"`
	LastActivity float64  `json:"last_activity"`
}
