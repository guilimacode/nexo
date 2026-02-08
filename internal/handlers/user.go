package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/models"
	"github.com/guilimacode/nexo/internal/store"
	"github.com/guilimacode/nexo/internal/utils"
)

func UpdateUserHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	currentUser, err := store.GetUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	if input.Email != currentUser.Email {
		exists, _ := store.CheckEmailUniqueForUpdate(input.Email, id)
		if exists {
			return c.Status(400).JSON(fiber.Map{"error": "Email já em uso por outro usuário"})
		}
	}

	currentUser.FullName = input.FullName
	currentUser.Role = input.Role
	currentUser.Email = input.Email

	if err := store.UpdateUser(currentUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar"})
	}

	return c.JSON(fiber.Map{"message": "Usuário atualizado", "user": currentUser})
}

func UpdatePasswordHandler(c *fiber.Ctx) error {
	idParam, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	loggedInUserID := c.Locals("userID").(float64)

	if int64(loggedInUserID) != idParam {
		return c.Status(403).JSON(fiber.Map{"error": "Você não tem permissão para alterar a senha de outro usuário"})
	}

	input := new(models.ChangePasswordDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	if len(input.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "A nova senha deve ter no mínimo 6 caracteres"})
	}

	user, err := store.GetUserById(idParam)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	if !utils.CheckPasswordHash(input.OldPassword, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "A senha antiga está incorreta"})
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao processar nova senha"})
	}

	if err := store.UpdatePassword(idParam, newHash); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar senha"})
	}

	return c.JSON(fiber.Map{"message": "Senha alterada com sucesso!"})
}

func DeleteUserHandler(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	if err := store.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar"})
	}

	return c.JSON(fiber.Map{"message": "Usuário deletado"})
}
