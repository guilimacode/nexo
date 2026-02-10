package models

import "time"

type User struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  int64     `db:"organization_id" json:"organization_id"`
	EstablishmentID *int64    `db:"establishment_id" json:"establishment_id"`
	FullName        string    `db:"full_name" json:"full_name"`
	Email           string    `db:"email" json:"email"`
	PasswordHash    string    `db:"password_hash" json:"-"`
	Role            string    `db:"role" json:"role"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
