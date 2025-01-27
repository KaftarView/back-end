package dto

import (
	"time"
)

type UserDetailsResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CouncilorsDetailsResponse struct {
	ID           uint   `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Profile      string `json:"profile"`
	EnteringYear int    `json:"enteringYear"`
	Description  string `json:"description"`
}

type RoleDetailsResponse struct {
	ID          uint            `json:"id"`
	CreatedAt   time.Time       `json:"createdAt"`
	Type        string          `json:"type"`
	Permissions map[uint]string `json:"permissions"`
}

type PermissionDetailsResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
