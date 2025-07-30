package main

import (
	"stream/internal/database"
	"stream/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()
	database.Connect()
	defer database.Client.Close()
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
