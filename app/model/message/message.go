package model

import "time"

type DirectMessage struct {
	ID         int64      `json:"id"`
	SenderID   int64      `json:"sender_id"`
	ReceiverID int64      `json:"receiver_id"`
	Message    string     `json:"message"`
	SentAt     time.Time  `json:"sent_at"`
	IsEdited   bool       `json:"is_edited"`
	EditedAt   *time.Time `json:"edited_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
