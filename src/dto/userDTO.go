package dto

import "time"

type UserDetailsResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
