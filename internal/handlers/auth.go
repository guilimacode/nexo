package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/models"
	"github.com/guilimacode/nexo/internal/store"
	"github.com/guilimacode/nexo/internal/utils"
)

func SignUpHandler(c *fiber.Ctx) error {
	dto := new(models.CreateUserDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	exists, _ := store.CheckEmailExists(dto.Email)
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "Este email já está em uso"})
	}

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao processar senha"})
	}

	user := models.User{
		FullName:       dto.Name + "" + dto.Surname,
		Email:          dto.Email,
		PasswordHash:   hash,
		Role:           dto.Role,
		OrganizationID: 1,
	}

	if err := store.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar usuário"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Cadastro Realizado!"})
}

func LoginHandler(c *fiber.Ctx) error {
	input := new(models.LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email e senha são obrigatórios"})
	}

	user, err := store.GetUserByEmail(input.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciais inválidas"})
	}

	match := utils.CheckPasswordHash(input.Password, user.PasswordHash)
	if !match {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciais inválidas"})
	}

	token, err := utils.GenerateToken(user.Email, user.Role, user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao gerar sessão"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(8 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Strict"

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Login realizado",
		"token":   token,
		"user":    user})
}
