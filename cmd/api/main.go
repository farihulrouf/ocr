package main

import (
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := configs.LoadConfig()

	// Init Dependencies
	configs.ConnectDB(cfg)
	configs.ConnectRedis(cfg)
	configs.InitS3(cfg)

	// Fiber Instance
	app := fiber.New(fiber.Config{
		AppName: "SEIDO OCR Enterprise API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Static files & Routes
	app.Static("/uploads", "./uploads")
	routes.SetupRoutes(app)

	log.Printf("🚀 SEIDO API started on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
