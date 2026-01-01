package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized, missing token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("tenant_id", claims["tenant_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

func SuperAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role != "ADMIN" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden: Super Admin only",
			})
		}

		return c.Next()
	}
}

func TenantAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role != "MANAGER" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden – Tenant Admin only",
			})
		}
		return c.Next()
	}
}

func EmployeeOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != "EMPLOYEE" {
			return c.Status(403).JSON(fiber.Map{"error": "Employees only"})
		}
		return c.Next()
	}
}
