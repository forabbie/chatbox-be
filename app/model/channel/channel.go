package model

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
