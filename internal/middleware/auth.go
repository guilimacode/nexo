package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guilimacode/nexo/internal/utils"
)

func Protected(c *fiber.Ctx) error {
	var tokenString string

	tokenString = c.Cookies("jwt")

	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		}
	}

	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Não autorizado: Token ausente"})
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Não autorizado: Token inválido"})
	}

	c.Locals("userID", claims["sub"])
	c.Locals("org", claims["org"])
	c.Locals("role", claims["role"])
	c.Locals("email", claims["email"])

	return c.Next()
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)

		if !ok {
			return c.Status(500).JSON(fiber.Map{"error": "Erro interno de permissão"})
		}

		for _, role := range allowedRoles {
			if role == userRole {
				return c.Next()
			}
		}
		return c.Status(401).JSON(fiber.Map{"error": "Acesso negado: Permissão insuficiente"})
	}
}
