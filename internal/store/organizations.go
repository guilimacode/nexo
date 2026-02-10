package store

import (
	"errors"

	"github.com/guilimacode/nexo/internal/models"
)

func CreateOrganization(o *models.Organization) error {
	query := `INSERT INTO organizations(name, document, created_at, updated_at) VALUES (:name, :document, NOW(), NOW()) RETURNING id`

	rows, err := DB.NamedQuery(query, o)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&o.ID); err != nil {
			return err
		}
	}

	return nil
}

func DeleteOrganization(id int64) error {
	query := `DELETE FROM organizations WHERE id = $1`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("organização não encontrada para deletar")
	}

	return nil
}
