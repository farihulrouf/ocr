package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized, missing token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Parsing Token (Secret key harusnya dari .env)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("SUPER_SECRET_KEY_JAPAN_2024"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Simpan data user & tenant ke context untuk dipakai di Handler
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("tenant_id", claims["tenant_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}