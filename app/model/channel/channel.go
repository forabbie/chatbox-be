package model

import "time"

// type InsertChannelParam struct {
// 	ID      int64   `json:"id"`
// 	Name    string  `json:"name"`
// 	UserIDs []int64 `json:"user_ids,omitempty"`
// }

type CreatePayload struct {
	Name    string  `json:"name" validate:"required"`
	UserIDs []int64 `json:"user_ids" validate:"required"`
}

type ChannelParam struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	CreatedBy int64   `json:"created_by_id"`
	UserIDs   []int64 `json:"user_ids"`
}

type GetChannelParam struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	CreatedBy     int64         `json:"created_by_id"`
	CreatedByUser UserSummary   `json:"created_by"`
	Members       []UserSummary `json:"members"`
}

type UserSummary struct {
	ID           int64  `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Emailaddress string `json:"emailaddress"`
}

type AddMemberRequest struct {
	ID       int64 `json:"id"`        // Channel ID
	MemberID int64 `json:"member_id"` // User ID to add as member
}

type ChannelWithMessage struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedBy int64      `json:"created_by"`
	UserIDs   []int64    `json:"user_ids"`
	MessageID *int64     `json:"message_id,omitempty"`
	Message   *string    `json:"message,omitempty"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
}
