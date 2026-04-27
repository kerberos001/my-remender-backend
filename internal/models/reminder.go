package models

import "time"

type Reminder struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	IsCompleted bool      `json:"isCompleted"`
	Priority    int32     `json:"priority"`
	CreatedAt   time.Time `json:"createdAt"`
}
