package teleapi

// Update ...
type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

// Message ...
type Message struct {
	MessageID int64  `json:"message_id"`
	From      User   `json:"from"`
	Date      int64  `json:"date"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

// User ...
type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat ...
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type getUpdatesResp struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}