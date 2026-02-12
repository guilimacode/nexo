package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/models"
	"github.com/guilimacode/nexo/internal/models/dtos"
	"github.com/guilimacode/nexo/internal/store"
	"github.com/guilimacode/nexo/internal/utils"
)

func CreateEstablishmentHandler(c *fiber.Ctx) error {
	orgIDVal := c.Locals("org")

	if orgIDVal == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido ou sem organização"})
	}
	orgID := int64(orgIDVal.(float64))

	dto := new(dtos.EstablishmentInputDTO)

	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	if dto.Document != "" {
		cleanDoc, err := utils.FormatAndValidateCpfCnpj(dto.Document)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Documento (CPF/CNPJ) inválido"})
		}
		dto.Document = cleanDoc
	}

	exists, _ := store.CheckDocumentExists(dto.Document)
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "CPF/CNPJ já registrado."})
	}

	address := &models.Address{
		Street:       dto.Address.Street,
		Complement:   dto.Address.Complement,
		Number:       dto.Address.Number,
		Neighborhood: dto.Address.Neighborhood,
		City:         dto.Address.City,
		State:        dto.Address.State,
		ZipCode:      dto.Address.ZipCode,
	}

	if err := store.CreateAddress(address); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar estabelecimento", "details": err.Error()})
	}

	estab := &models.Establishment{
		OrganizationID: orgID,
		AddressID:      &address.ID,
		Name:           dto.Name,
		Nickname:       dto.Nickname,
		Document:       dto.Document,
	}

	if err := store.CreateEstablishment(estab); err != nil {
		store.DeleteAddress(address.ID)
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar estabelecimento", "details": err.Error()})
	}

	estab.Address = address

	return c.Status(201).JSON(fiber.Map{"message": "Estabelecimento criado com sucesso!", "estabelecimento": estab})
}

func ListEstablishmentsHandler(c *fiber.Ctx) error {
	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	orgID := int64(orgIDVal)

	list, err := store.GetEstablishmentsByOrg(orgID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao listar estabelecimentos"})
	}

	return c.JSON(list)
}

func GetEstablishmentsByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	requesterOrgID := int64(orgIDVal)

	establishment, err := store.GetEstablismentById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Estabelecimento nâo encontrado"})
	}

	if establishment.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Esse estabelecimento nâo pertence a sua organização"})
	}

	if establishment.AddressID != nil {
		address, _ := store.GetAddressById(*establishment.AddressID)
		establishment.Address = address
	}

	return c.JSON(establishment)
}

func UpdateEstablishmentHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	requesterOrgID := int64(orgIDVal)

	currentEstablishment, err := store.GetEstablismentById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Estabelecimento nâo encontrado"})
	}

	if currentEstablishment.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Esse estabelecimento nâo pertence a sua organização"})
	}

	dto := new(dtos.EstablishmentInputDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos"})
	}

	if dto.Document != "" {
		cleanDoc, err := utils.FormatAndValidateCpfCnpj(dto.Document)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Documento inválido"})
		}
		dto.Document = cleanDoc
		if dto.Document != currentEstablishment.Document {
			exists, _ := store.CheckDocumentForUpdate(dto.Document, id)
			if exists {
				return c.Status(409).JSON(fiber.Map{"error": "Este CPF/CNPJ já está em uso por outro estabelecimento."})
			}
		}
	}

	address := &models.Address{
		ID:           *currentEstablishment.AddressID,
		Street:       dto.Address.Street,
		Complement:   dto.Address.Complement,
		Number:       dto.Address.Number,
		Neighborhood: dto.Address.Neighborhood,
		City:         dto.Address.City,
		State:        dto.Address.State,
		ZipCode:      dto.Address.ZipCode,
	}

	if err := store.UpdateAddress(address); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar endereço"})
	}

	currentEstablishment.Name = dto.Name
	currentEstablishment.Nickname = dto.Nickname
	currentEstablishment.Document = dto.Document

	if err := store.UpdateEstablishment(currentEstablishment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar estabelecimento"})
	}

	currentEstablishment.Address = address
	return c.JSON(fiber.Map{"message": "Atualizado com sucesso", "estabelecimento": currentEstablishment})
}

func DeleteEstablishmentHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	orgIDVal, ok := c.Locals("org").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Token inválido"})
	}
	requesterOrgID := int64(orgIDVal)

	establishment, err := store.GetEstablismentById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Estabelecimento nâo encontrado"})
	}

	if establishment.OrganizationID != requesterOrgID {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso negado: Esse estabelecimento nâo pertence a sua organização"})
	}

	if establishment.AddressID != nil {
		address, _ := store.GetAddressById(*establishment.AddressID)
		establishment.Address = address
	}

	if err := store.DeleteEstablishment(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar estabelecimento"})
	}

	return c.JSON(fiber.Map{"message": "Estabelecimento deletado", "estabelecimento": establishment})
}
