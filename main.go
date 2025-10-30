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
		AppName:      "Thai Location API v1.0.0",
		ServerHeader: "Thai Location API",
		// Trust proxy headers (for Cloudflare)
		EnableTrustedProxyCheck: true,
		TrustedProxies: []string{
			"173.245.48.0/20",
			"103.21.244.0/22",
			"103.22.200.0/22",
			"103.31.4.0/22",
			"141.101.64.0/18",
			"108.162.192.0/18",
			"190.93.240.0/20",
			"188.114.96.0/20",
			"197.234.240.0/22",
			"198.41.128.0/17",
			"162.158.0.0/15",
			"104.16.0.0/13",
			"104.24.0.0/14",
			"172.64.0.0/13",
			"131.0.72.0/22",
		},
		ProxyHeader: "CF-Connecting-IP,X-Forwarded-For",
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
	}))
	app.Use(recover.New())

	// Enhanced CORS for Cloudflare
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "*",
		AllowCredentials: false,
		ExposeHeaders:    "X-Request-ID,X-Response-Time",
		MaxAge:           86400, // 24 hours
	}))

	// Security headers middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Cache control for API responses
		if c.Path() != "/health" {
			c.Set("Cache-Control", "public, max-age=300, s-maxage=3600")
		}

		return c.Next()
	})

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
			"status":  "OK",
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

	// Serve API documentation (Swagger UI)
	// Static files are located in ./docs (index.html + openapi.json)
	app.Static("/docs", "./docs")

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
