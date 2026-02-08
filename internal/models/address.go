package models

import "time"

type Address struct {
	ID         int64     `db:"id" json:"id"`
	Street     string    `db:"street" json:"street"`
	Complement string    `db:"complement" json:"complement"`
	Number     int       `db:"number" json:"number"`
	District   string    `db:"district" json:"district"`
	City       string    `db:"city" json:"city"`
	State      string    `db:"state" json:"state"`
	ZipCode    string    `db:"zip_code" json:"zip_code"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CreateAddressDTO struct {
	Street     string `json:"street" validate:"required"`
	Complement string `json:"complement"`
	Number     int    `json:"number" validate:"required,numeric"`
	District   string `json:"district" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required,len=2,uppercase"`
	ZipCode    string `json:"zip_code" validate:"required,len=8,numeric"`
}
