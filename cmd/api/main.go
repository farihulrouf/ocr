package main

import (
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load .env
	configs.LoadConfig()

	// DB
	configs.ConnectDB()

	configs.ConnectRedis() // <==== wajib ini sebelum router
	app := fiber.New()

	// ✅ CORS ALLOW ALL (DEV ONLY)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	// Routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
