package store

import (
	"errors"

	"github.com/guilimacode/nexo/internal/models"
)

func CreateAddress(a *models.Address) error {
	query := `INSERT INTO addresses(street, complement, number, district, city, state, zip_code, created_at, updated_at) VALUES (:street, :complement, :number, :district, :city, :state, :zip_code, NOW(), NOW()) RETURNING id`

	rows, err := DB.NamedQuery(query, a)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&a.ID); err != nil {
			return err
		}
	}
	return nil
}

func UpdateAddress(a *models.Address) error {
	query := `
			UPDATE ADDRESSES
			SET street = :street,
				complement = :complement,
				number = :number,
				district = :district,
				city = :city,
				state = :state,
				zip_code = :zip_code,
				updated_at = NOW()
			WHERE id = :id`

	_, err := DB.NamedExec(query, a)

	return err
}

func GetAddressById(id int64) (*models.Address, error) {
	var address models.Address
	query := `SELECT * FROM addresses WHERE id = $1`

	err := DB.Get(&address, query, id)

	if err != nil {
		return nil, errors.New("endereço não encontrado")
	}

	return &address, nil
}

func DeleteAddress(id int64) error {
	_, err := DB.Exec("DELETE FROM addresses WHERE id = $1", id)
	return err
}
