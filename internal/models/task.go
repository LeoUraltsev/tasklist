package models

import (
	"errors"
	"time"
)

type Status string

var (
	Pending Status = "Pending"
	Done    Status = "Done"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
