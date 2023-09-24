package model

import "time"

type User struct {
	ID         int64     `json:"user_id,omitempty"`
	Fullname   string    `json:"fullname" binding:"required"`
	Email      string    `json:"email" binding:"required,email"`
	Password   string    `json:"password,omitempty" binding:"required"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
