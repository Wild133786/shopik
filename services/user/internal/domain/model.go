package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type RegisterInput struct {
	Email    string
	Password string
}
type LoginInput struct {
	Email    string
	Password string
}
