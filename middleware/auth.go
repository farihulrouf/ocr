package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected: Validasi Token Dasar & Ekstrak Claims
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

// FinanceOnly: Role Tertinggi (Bisa akses segalanya termasuk fitur Finance)
func FinanceOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		// Tambahin MANAGER di sini
		if role == "FINANCE" || role == "ADMIN" || role == "MANAGER" {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden: Finance/Manager access required",
		})
	}
}

// TenantAdminOnly: Untuk Manager, tapi Finance & Admin juga BOLEH tembus
func TenantAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		// Karena Finance paling tinggi, mereka harus bisa akses menu Manager
		if role == "MANAGER" || role == "FINANCE" || role == "ADMIN" {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden: Manager or Higher role required",
		})
	}
}

// SuperAdminOnly: Khusus untuk maintenance sistem (System Level)
func SuperAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		if role != "ADMIN" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden: System Admin only",
			})
		}

		return c.Next()
	}
}

// EmployeeOnly: Role paling dasar (Hanya bisa akses data milik sendiri)
func EmployeeOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		// Employee akses terbatas, tapi Finance/Manager biasanya boleh intip buat monitoring
		if role == "EMPLOYEE" || role == "MANAGER" || role == "FINANCE" || role == "ADMIN" {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}
}
