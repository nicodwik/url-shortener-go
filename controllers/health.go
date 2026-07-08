package controllers

import (
	"net/http"
	"url-shortener-go/config"
	"url-shortener-go/helpers"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Cache    string `json:"cache"`
	Message  string `json:"message"`
}

func HealthCheck(c echo.Context) error {
	healthResponse := HealthResponse{
		Status:   "healthy",
		Database: "connected",
		Cache:    "connected",
		Message:  "Service is running properly",
	}

	// Verify database connection
	db, err := config.GetDB()
	if err != nil || db == nil {
		healthResponse.Database = "disconnected"
		healthResponse.Status = "unhealthy"
	}

	// Verify cache connection
	cache, err := config.GetCache()
	if err != nil || cache == nil {
		healthResponse.Cache = "disconnected"
		healthResponse.Status = "unhealthy"
	}

	statusCode := http.StatusOK
	if healthResponse.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, helpers.ResponseOK("Health check", healthResponse))
}
