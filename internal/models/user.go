package models

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
}
