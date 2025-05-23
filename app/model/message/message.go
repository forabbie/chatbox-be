package model

import "time"

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	EmailAddress string `json:"emailaddress"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
}

type Message struct {
	ID            int64      `json:"id,omitempty"`
	Message       string     `json:"message"`
	SentAt        time.Time  `json:"sent_at"`
	IsEdited      bool       `json:"is_edited"`
	EditedAt      *time.Time `json:"edited_at,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Sender        User       `json:"sender"`
	Receiver      *User      `json:"receiver,omitempty"`
	ReceiverID    *int64     `json:"receiver_id"`
	ReceiverClass string     `json:"receiver_class"`
}

type Query struct {
	// Message   *string `json:"message,omitempty" query:"message"`
	// Firstname *string `json:"firstname,omitempty" query:"firstname"`
	// Lastname  *string `json:"lastname,omitempty" query:"lastname"`
	// Username  *string `json:"username,omitempty" query:"username"`

	// SenderID   *string `json:"sender_id,omitempty" query:"sender_id"`
	ReceiverID    *int64  `json:"receiver_id,omitempty" query:"receiver_id"`
	ReceiverClass *string `query:"receiver_class"`
	Created       struct {
		Gte *time.Time `json:"gte,omitempty" query:"gte"`
		Lte *time.Time `json:"lte,omitempty" query:"lte"`
	} `json:"created,omitempty" query:"created"`
}
