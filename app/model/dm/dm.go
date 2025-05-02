package model

import "time"

type DMListItem struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Message    string    `json:"message"`
	SentAt     time.Time `json:"sent_at"`
}
