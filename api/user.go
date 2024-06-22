package api

import "time"

type createUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password"  binding:"required,min=6"`
	FullName string `json:"full_name"  binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
}

type userResponse struct {
	UserName          string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name"  binding:"required"`
	Email             string    `json:"email"  binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
