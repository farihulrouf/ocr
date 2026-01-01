package main

import (
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load .env
	configs.LoadConfig()

	// 2. Koneksi DB + Migrasi
	configs.ConnectDB()

	// 3. Jalankan Seeder
	//configs.SeedDatabase(configs.DB)

	app := fiber.New()

	// 4. Routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
