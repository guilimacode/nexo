package store

import (
	"errors"

	"github.com/guilimacode/nexo/internal/models"
)

func CreateUser(u *models.User) error {
	query := `INSERT INTO users(organization_id, establishment_id, full_name, email, password_hash, role, created_at, updated_at) VALUES (:organization_id, :establishment_id, :full_name, :email, :password_hash, :role, NOW(), NOW()) RETURNING id`

	rows, err := DB.NamedQuery(query, u)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&u.ID); err != nil {
			return err
		}
	}
	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	err := DB.Get(&user, query, email)

	if err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	return &user, nil
}

func GetAllUsers(organization_id int64) ([]models.User, error) {
	users := []models.User{}

	query := `SELECT * FROM users WHERE organization_id = $1 ORDER BY full_name ASC`

	err := DB.Select(&users, query, organization_id)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserById(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := DB.Get(&user, query, id)

	if err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	return &user, nil
}

func UpdateUser(u *models.User) error {
	query := `
				UPDATE users
				SET full_name = :full_name,
				role = :role,
				establishment_id = :establishment_id,
				updated_at = NOW()
				WHERE id = :id`

	_, err := DB.NamedExec(query, u)

	return err
}

func UpdatePassword(id int64, newHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := DB.Exec(query, newHash, id)
	return err
}

func DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("usuário não encontrado para deletar")
	}

	return nil
}

func CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 from users WHERE email = $1)`
	err := DB.Get(&exists, query, email)

	return exists, err
}

func CheckEmailUniqueForUpdate(email string, excludeId int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 from users WHERE email = $1 AND id != $2)`
	err := DB.Get(&exists, query, email, excludeId)
	return exists, err
}
