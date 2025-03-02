package models

import "time"

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
