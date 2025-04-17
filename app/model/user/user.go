package model

import "time"

type User struct {
	Id           int        `json:"id,omitempty"`
	Firstname    string     `json:"firstname,omitempty" validate:"required"`
	Lastname     string     `json:"lastname,omitempty" validate:"required"`
	EmailAddress string     `json:"emailaddress,omitempty" validate:"required,emailaddress"`
	Username     string     `json:"username,omitempty"`
	Password     string     `json:"hashed_password,omitempty"`
	IsActive     *bool      `json:"is_active,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

type Query struct {
	Id           int     `json:"id,omitempty"`
	Firstname    *string `json:"firstname,omitempty" query:"firstname"`
	Lastname     *string `json:"lastname,omitempty" query:"lastname"`
	EmailAddress *string `json:"emailaddress,omitempty" query:"emailaddress"`
	Username     *string `json:"username,omitempty" query:"username"`
	IsActive     *bool   `json:"is_active,omitempty" query:"is_active"`
	Created      struct {
		Gte *time.Time `json:"gte,omitempty" query:"gte"`
		Lte *time.Time `json:"lte,omitempty" query:"lte"`
	} `json:"created,omitempty" query:"created"`
}

type UserDetails struct {
	Id           int        `json:"id,omitempty"`
	Firstname    string     `json:"firstname,omitempty" validate:"required"`
	Lastname     string     `json:"lastname,omitempty" validate:"required"`
	EmailAddress string     `json:"emailaddress,omitempty" validate:"required,emailaddress"`
	Username     string     `json:"username,omitempty"`
	IsActive     *bool      `json:"is_active,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}
