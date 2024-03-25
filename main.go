package main

import (
	"log"
	"url-shortener-go/config"
	"url-shortener-go/controllers"
	dashboardRepo "url-shortener-go/repository"
	"url-shortener-go/routes"
	"url-shortener-go/scheduler"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitCache()
	config.InitDB()

	go scheduler.Init()

	// dependency injection
	dashboardRepo := dashboardRepo.InitDasboardRepository(config.DBConn)
	dashboardController := controllers.NewDashboardController(dashboardRepo)

	e := echo.New()
	routes.Init(e, *dashboardController)

	if err := e.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}
