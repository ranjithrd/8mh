package main

import (
	"backend/src/db"
	_ "backend/src/docs"
	"backend/src/handlers"
	"backend/src/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title 8MH API
// @version 1.0
// @description API for 8MH cooperative banking system with blockchain verification
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@8mh.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name session_id

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	handlers.InitAuthHandlers()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	routes.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	e.Logger.Fatal(e.Start(":" + port))
}
