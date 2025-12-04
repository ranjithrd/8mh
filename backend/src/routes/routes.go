package routes

import (
	"backend/src/handlers"
	"backend/src/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}

type IndexResponse struct {
	Message string `json:"message" example:"API is running"`
	Version string `json:"version" example:"1.0.0"`
}

// Index godoc
// @Summary API Index
// @Description Returns API status and version
// @Tags general
// @Produce json
// @Success 200 {object} IndexResponse
// @Router / [get]
func index(c echo.Context) error {
	return c.JSON(http.StatusOK, IndexResponse{
		Message: "API is running",
		Version: "1.0.0",
	})
}

// Health godoc
// @Summary Health Check
// @Description Returns health status of the API
// @Tags general
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func health(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: "healthy",
	})
}

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.RateLimiter(100))

	e.GET("/", index)
	e.GET("/health", health)

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/login", handlers.Login, middleware.RateLimiter(8))
	auth.POST("/verify", handlers.Verify)
	auth.GET("/me", handlers.Me, middleware.Auth)
}
