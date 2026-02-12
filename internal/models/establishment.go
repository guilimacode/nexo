package models

import "time"

type Establishment struct {
	ID             int64     `db:"id"              json:"id"`
	OrganizationID int64     `db:"organization_id" json:"organization_id"`
	AddressID      *int64    `db:"address_id"      json:"address_id"`
	Name           string    `db:"name"            json:"name"`
	Nickname       string    `db:"nickname"        json:"nickname"`
	Document       string    `db:"document"        json:"document"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"      json:"updated_at"`
	Address        *Address  `db:"-"               json:"address,omitempty"`
}
