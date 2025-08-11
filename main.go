package main

import (
	"fmt"
	"log"
	"url-shortener-go/config"
	"url-shortener-go/controllers"
	"url-shortener-go/helpers"
	dashboardRepo "url-shortener-go/repository"
	"url-shortener-go/routes"
	"url-shortener-go/scheduler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := config.InitCache(); err != nil {
		log.Fatalf("Error connecting to cache: %v", err)
	}

	if _, err := config.InitDB(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	go scheduler.Init()

	// dependency injection
	dashboardRepo := dashboardRepo.InitDasboardRepository(config.DBConn)
	dashboardController := controllers.NewDashboardController(dashboardRepo)

	e := echo.New()
	routes.Init(e, *dashboardController)

	appPort := fmt.Sprintf(":%v", helpers.Env("APP_PORT", "8000"))
	if err := e.Start(appPort); err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
