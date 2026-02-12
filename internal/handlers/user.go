package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/models"
	"github.com/guilimacode/nexo/internal/models/dtos"
	"github.com/guilimacode/nexo/internal/store"
	"github.com/guilimacode/nexo/internal/utils"
)

func CreateUserHandler(c *fiber.Ctx) error {
	establishmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	requesterOrgID := int64(orgIDVal)
	requesterID := int64(c.Locals("sub").(float64))
	requesterRole := c.Locals("role").(string)

	establishment, err := store.GetEstablismentById(establishmentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Estabelecimento não encontrado"})
	}

	if establishment.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Estabelecimento de outra organização"})
	}

	if requesterRole != "owner" {
		requester, err := store.GetUserById(requesterID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro interno"})
		}
		if requester.EstablishmentID == nil || *requester.EstablishmentID != establishmentID {
			return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Você não gerencia este estabelecimento"})
		}
	}

	dto := new(dtos.CreateUserDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	if requesterRole == "manager" {
		if dto.Role == "owner" || dto.Role == "manager" {
			return c.Status(403).JSON(fiber.Map{"error": "Gerentes só podem criar colaboradores/vendedores"})
		}
	}

	emailExists, _ := store.CheckEmailExists(dto.Email)
	if emailExists {
		return c.Status(409).JSON(fiber.Map{"error": "Email já cadastrado"})
	}

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao processar senha"})
	}

	user := &models.User{
		OrganizationID:  requesterOrgID,
		EstablishmentID: &establishmentID,
		FullName:        dto.FullName,
		Email:           dto.Email,
		PasswordHash:    hash,
		Role:            dto.Role,
	}

	if err := store.CreateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar usuário"})
	}

	return c.Status(201).JSON(user)
}

func ListUsersHandler(c *fiber.Ctx) error {
	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	orgID := int64(orgIDVal)
	requesterRole := c.Locals("role").(string)

	if requesterRole != "owner" {
		return c.Status(403).JSON(fiber.Map{"error": "Permissão insuficiente"})
	}

	list, err := store.GetUsersByOrg(orgID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao listar usuários"})
	}

	return c.JSON(list)
}

func ListUsersByEstablishmentHandler(c *fiber.Ctx) error {
	establishmentID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	requesterOrgID := int64(orgIDVal)
	requesterID := int64(c.Locals("sub").(float64))
	requesterRole := c.Locals("role").(string)

	establishment, err := store.GetEstablismentById(establishmentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Esse estabelecimento não existe"})
	}

	if establishment.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Estabelecimento de outra organização"})
	}

	if requesterRole != "owner" {
		requester, err := store.GetUserById(requesterID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao validar permissões"})
		}
		if requester.EstablishmentID == nil || *requester.EstablishmentID != establishmentID {
			return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Você não pertence a este estabelecimento"})
		}
	}

	list, err := store.GetUsersByEstablishment(establishmentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao listar usuários"})
	}

	return c.JSON(list)
}

func UpdateUserHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	requesterOrgID := int64(c.Locals("org").(float64))

	targetUser, err := store.GetUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	if targetUser.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Usuário pertence a outra organização"})
	}

	dto := new(dtos.UpdateUserDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	if dto.Email != targetUser.Email {
		exists, _ := store.CheckEmailUniqueForUpdate(dto.Email, id)
		if exists {
			return c.Status(400).JSON(fiber.Map{"error": "Email já em uso por outro usuário"})
		}
	}

	targetUser.FullName = dto.FullName
	targetUser.Role = dto.Role
	targetUser.Email = dto.Email

	if err := store.UpdateUser(targetUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar"})
	}

	return c.JSON(fiber.Map{"message": "Usuário atualizado", "user": targetUser})
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

	input := new(dtos.ChangePasswordDTO)
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

	requesterOrgID := int64(c.Locals("org").(float64))
	requesterRole := c.Locals("role").(string)

	targetUser, err := store.GetUserById(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuário não encontrado"})
	}

	if targetUser.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado"})
	}

	if requesterRole == "seller" {
		return c.Status(403).JSON(fiber.Map{"error": "Permissão insuficiente"})
	}

	requesterID := int64(c.Locals("sub").(float64))
	if targetUser.ID == requesterID {
		return c.Status(400).JSON(fiber.Map{"error": "Você não pode deletar seu próprio usuário"})
	}

	if err := store.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar"})
	}

	return c.JSON(fiber.Map{"message": "Usuário deletado"})
}
