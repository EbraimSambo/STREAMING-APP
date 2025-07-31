package main

import (
	"log"
	"stream/internal/config"
	"stream/internal/database"
	"stream/internal/routes"
	"stream/internal/worker"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.Connect()
	defer database.Client.Close()

	// Start workers
	worker.StartDispatcher(5, database.Client, cfg) // 5 workers, adjust as needed

	e := echo.New()

	// Habilita o CORS para todas as rotas
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://seusite.com"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	routes.Routes(e, database.Client)
	e.Logger.Fatal(e.Start(":3344"))
}