package model

type Channel struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	UserIDs []int64 `json:"user_ids,omitempty"`
}

type CreatePayload struct {
	CreatedBy int64   `json:"created_by,omitempty"`
	Name      string  `json:"name" validate:"required"`
	UserIDs   []int64 `json:"user_ids" validate:"required"`
}
