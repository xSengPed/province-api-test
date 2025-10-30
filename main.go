package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Thai Location API v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	// Initialize data service
	dataService, err := NewDataService("./data/raw")
	if err != nil {
		log.Fatal("Failed to initialize data service:", err)
	}

	// Initialize handlers
	handler := NewLocationHandler(dataService)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
			"message": "Thai Location API is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	
	// Geography routes
	api.Get("/geographies", handler.GetGeographies)
	
	// Province routes
	api.Get("/provinces", handler.GetProvinces)
	api.Get("/provinces/:id", handler.GetProvinceByID)
	api.Get("/provinces/:id/districts", handler.GetDistrictsByProvinceID)
	
	// District routes
	api.Get("/districts", handler.GetDistricts)
	api.Get("/districts/:id", handler.GetDistrictByID)
	api.Get("/districts/:id/subdistricts", handler.GetSubDistrictsByDistrictID)
	
	// Sub-district (Tambon) routes
	api.Get("/subdistricts", handler.GetSubDistricts)
	api.Get("/subdistricts/:id", handler.GetSubDistrictByID)

	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}