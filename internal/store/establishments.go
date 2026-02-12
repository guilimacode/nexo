package store

import (
	"errors"

	"github.com/guilimacode/nexo/internal/models"
)

func CreateEstablishment(e *models.Establishment) error {
	query := `INSERT INTO establishments(organization_id, address_id, name, nickname, document, created_at, updated_at) VALUES(:organization_id, :address_id, :name, :nickname, :document, NOW(), NOW()) RETURNING id`

	rows, err := DB.NamedQuery(query, e)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&e.ID); err != nil {
			return err
		}
	}
	return nil
}

func GetEstablishmentsByOrg(organization_id int64) ([]models.Establishment, error) {
	establishments := []models.Establishment{}

	query := `SELECT * FROM establishments WHERE organization_id = $1`

	err := DB.Select(&establishments, query, organization_id)
	if err != nil {
		return nil, err
	}

	return establishments, nil
}

func GetEstablismentById(id int64) (*models.Establishment, error) {
	var establishment models.Establishment

	query := `SELECT * FROM establishments WHERE id = $1`

	err := DB.Get(&establishment, query, id)

	if err != nil {
		return nil, errors.New("estabelecimento não encontrado")
	}

	return &establishment, nil
}

func UpdateEstablishment(e *models.Establishment) error {
	query := `
	UPDATE establishments
		SET name = :name,
			nickname = :nickname,
			document = :document,
			updated_at = NOW()
		WHERE id = :id AND organization_id = :organization_id`

	_, err := DB.NamedExec(query, e)
	return err
}

func DeleteEstablishment(id int64) error {
	var addressID int64
	queryGet := `SELECT address_id FROM establishments WHERE id = $1`
	_ = DB.Get(&addressID, queryGet, id)

	queryDelEst := `DELETE FROM establishments WHERE id = $1`
	_, err := DB.Exec(queryDelEst, id)
	if err != nil {
		return err
	}

	if addressID != 0 {
		queryDelAddr := `DELETE FROM addresses WHERE id = $1`
		_, err = DB.Exec(queryDelAddr, addressID)
		return err
	}

	return nil
}

func CheckDocumentExists(document string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM establishments WHERE document = $1)`
	err := DB.Get(&exists, query, document)

	return exists, err
}

func CheckDocumentForUpdate(document string, excludeId int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM establishments WHERE document = $1 AND id <> $2)`
	err := DB.Get(&exists, query, document, excludeId)

	return exists, err
}
