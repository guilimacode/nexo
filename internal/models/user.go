package models

import "time"

type User struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	EstablishmentID *int64    `db:"establishment_id" json:"establishment_id"`
	FullName        string    `db:"full_name" json:"full_name"`
	Email           string    `db:"email" json:"email"`
	PasswordHash    string    `db:"password_hash" json:"password_hash"`
	Role            string    `db:"role" json:"role"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserDTO struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Role            string `json:"role"`
	EstablishmentID int64  `json:"establishment_id"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
