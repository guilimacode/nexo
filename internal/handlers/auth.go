package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/models"
	"github.com/guilimacode/nexo/internal/models/dtos"
	"github.com/guilimacode/nexo/internal/store"
	"github.com/guilimacode/nexo/internal/utils"
)

var validate = validator.New()

func SignUpHandler(c *fiber.Ctx) error {
	dto := new(dtos.SignupDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
	}

	if dto.CompanyDocument != "" {
		cleanDoc, err := utils.FormatAndValidateCpfCnpj(dto.CompanyDocument)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Documento (CPF/CNPJ) inválido"})
		}
		dto.CompanyDocument = cleanDoc
	}

	exists, _ := store.CheckEmailExists(dto.Email)
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "Este email já está em uso"})
	}

	org := &models.Organization{
		Name:     dto.CompanyName,
		Document: dto.CompanyDocument,
	}

	if err := store.CreateOrganization(org); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar empresa", "details": err.Error()})
	}

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao processar senha"})
	}

	user := &models.User{
		FullName:        dto.FullName,
		Email:           dto.Email,
		PasswordHash:    hash,
		Role:            "owner",
		OrganizationID:  org.ID,
		EstablishmentID: nil,
	}

	if err := store.CreateUser(user); err != nil {
		store.DeleteOrganization(org.ID)
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar usuário"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":         "Conta criada com sucesso!",
		"user_id":         user.ID,
		"organization_id": org.ID,
	})
}

func LoginHandler(c *fiber.Ctx) error {
	input := new(dtos.LoginInput)
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

	token, err := utils.GenerateToken(user.Email, user.Role, user.ID, user.OrganizationID)
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
