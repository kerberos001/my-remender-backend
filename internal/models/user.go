// internal/models/user.go
package models

import "time"

// User representa la tabla 'users'
type User struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"` // Excluir del JSON
	BiometricsEnabled bool      `json:"biometricsEnabled"`
	CreatedAt         time.Time `json:"createdAt"`
}

// AuthPayload representa el tipo AuthPayload de GraphQL
type AuthPayload struct {
	Token     string `json:"token"`
	User      *User  `json:"user"`
	Status    string `json:"status"` // PENDING, APPROVED
	GroupName string `json:"groupName,omitempty"`
}
