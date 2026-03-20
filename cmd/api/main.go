package main

import (
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// ✅ Load config
	cfg := configs.LoadConfig()

	// ✅ Init dependencies (pakai cfg)
	configs.InitS3(cfg)
	configs.ConnectDB(cfg)
	configs.ConnectRedis(cfg)

	// ✅ Init Fiber
	app := fiber.New()

	// ✅ CORS (dev only)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	// ✅ Static file
	app.Static("/uploads", "./uploads")

	// ✅ Routes
	routes.SetupRoutes(app)

	log.Println("🚀 API running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
