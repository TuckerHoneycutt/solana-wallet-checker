package main

import (
	"log"
	"solana-wallet-checker/handlers"
	"solana-wallet-checker/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load bluechip token configuration
	configService, err := services.NewConfigService("config/bluechip_tokens.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers with config service
	h := handlers.NewHandlers(configService)

	// Routes
	e.GET("/", h.HomeHandler)
	e.GET("/balance", h.BalanceHandler)

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Loaded bluechip tokens:", len(configService.GetAllBluechipTokens()))
	log.Fatal(e.Start(":8080"))
}
