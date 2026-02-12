package models

import "time"

type Address struct {
	ID           int64     `db:"id" json:"id"`
	Street       string    `db:"street" json:"street"`
	Complement   string    `db:"complement" json:"complement"`
	Number       int       `db:"number" json:"number"`
	Neighborhood string    `db:"neighborhood" json:"neighborhood"`
	City         string    `db:"city" json:"city"`
	State        string    `db:"state" json:"state"`
	ZipCode      string    `db:"zip_code" json:"zip_code"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
