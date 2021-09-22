package entity

import (
	"time"
)

// User represents user model
type User struct {
	ID             string
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
